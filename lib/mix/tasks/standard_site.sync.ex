defmodule Mix.Tasks.StandardSite.Sync do
  use Mix.Task

  @shortdoc "Fetch Standard.site document records into content/standard_site_documents.yml"

  alias Site.StandardSite
  alias Site.StandardSiteDocuments

  @requirements ["app.start"]

  @impl Mix.Task
  def run(args) do
    {opts, _args, _invalid} =
      OptionParser.parse(args,
        strict: [
          did: :string,
          publication: :string,
          pds: :string,
          output: :string
        ]
      )

    did = opts[:did] || System.get_env("STANDARD_SITE_DID") || StandardSite.did()

    publication =
      opts[:publication] || System.get_env("STANDARD_SITE_PUBLICATION") ||
        StandardSite.publication()

    pds = opts[:pds] || System.get_env("STANDARD_SITE_PDS")
    output = opts[:output] || StandardSite.documents_path()

    records =
      did
      |> fetch_documents(publication, pds)
      |> localize_record_paths()

    documents =
      records
      |> StandardSiteDocuments.merge_records(StandardSiteDocuments.all(output))
      |> keep_local_posts_only()

    StandardSiteDocuments.write!(documents, output)

    Mix.shell().info("Wrote #{map_size(documents)} Standard.site document mappings to #{output}")
  end

  def fetch_documents(did, publication, pds \\ nil) do
    [did: did, pds: pds, collection: StandardSite.collection()]
    |> Site.AtProto.list_records()
    |> Enum.filter(fn %{"value" => value} ->
      publication == nil or Site.StandardSite.publication_uri(value) == publication
    end)
  end

  def localize_record_paths(records) do
    posts = Site.Blog.all_posts()
    by_path = Map.new(posts, fn post -> {post.path, post.path} end)
    by_legacy_path = Map.new(posts, fn post -> {"/" <> post.id, post.path} end)
    by_title = Map.new(posts, fn post -> {post.title, post.path} end)

    Enum.map(records, fn %{"value" => value} = record ->
      local_path =
        by_path[value["path"]] || by_legacy_path[value["path"]] || by_title[value["title"]] ||
          value["path"]

      put_in(record, ["value", "path"], local_path)
    end)
  end

  defp keep_local_posts_only(documents) do
    local_paths = MapSet.new(Enum.map(Site.Blog.all_posts(), & &1.path))

    Map.filter(documents, fn {path, _attrs} -> MapSet.member?(local_paths, path) end)
  end
end
