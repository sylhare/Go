# Hugo

Hugo is a go framework for building websites, a bit like Jekyll.

## Get Started with Hugo

Install [Hugo](https://gohugo.io/getting-started/quick-start/) with brew on mac:
```bash
brew install hugo
...
hugo version
>> Hugo Static Site Generator v0.51
```

And then create your first website with:
```bash
hugo new site quickstart
```

## Configure your site

Once you have created the project follow the instructions:

1. Download a theme into the same-named folder.
   Choose a [theme](https://themes.gohugo.io/), or
   create your own with the `hugo new theme <THEMENAME>` command.
2. You can add content by creating a single files
   with `hugo new <SECTIONNAME>/<FILENAME>.<FORMAT>`.
3. Start the built-in live server via `hugo server` and it will run at [localhost:1313](http://localhost:1313/)

## Sources

Here are a couple of useful links:

- [Hugo a static site generator](https://gohugo.io/)