export type Page = {
  slug: string;
  title: string;
};

export type WritingW = {
  data?: Writing;
};

export type Writing = {
  metadata: WritingMetadata;
  content: string;
};

export type WritingMetadata = {
  authors: string[];
  categories?: string[];
  date: string;
  read_time: number;
  slug: string;
  tags?: string[];
  title: string;
  type: string;
  bluesky_link?: string;
};

export type Project = {
  title: string;
  description: string;
  shortDescription?: string;
  link: string;
  image?: string;
  time?: string;
  languages?: string;
};

export type ProjectCategory = {
  category: string;
  items: Project[];
};
