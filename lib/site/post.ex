defmodule Site.Post do
  @enforce_keys [:id, :title, :body, :tags, :date, :path, :read_time]
  defstruct [:id, :title, :body, :tags, :date, :path, :read_time, :author, :description]

  def build(filename, attrs, body) do
    path = Path.rootname(filename)
    [year, month_day_id] = path |> Path.split() |> Enum.take(-2)
    [month, day, id] = String.split(month_day_id, "-", parts: 3)
    path = "/blog/" <> id
    date = Date.from_iso8601!("#{year}-#{month}-#{day}")

    struct!(
      __MODULE__,
      [
        id: id,
        date: date,
        body: body,
        path: path,
        read_time: read_time(body)
      ] ++ Map.to_list(attrs)
    )
  end

  defp read_time(body) do
    body
    |> String.replace(~r/<[^>]*>/, " ")
    |> String.split(~r/\s+/, trim: true)
    |> length()
    |> Kernel./(200)
    |> Float.ceil()
    |> trunc()
    |> max(1)
  end
end
