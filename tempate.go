package main

const tempate_html = `
<!DOCTYPE html>
<html lang="zh_CN" data-theme-mode="light" data-light-theme="{{ .Theme }}" data-dark-theme="midnight">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=0"/>
    <meta name="apple-mobile-web-app-capable" content="yes">
    <meta name="mobile-web-app-capable" content="yes"/>
    <meta name="apple-mobile-web-app-status-bar-style" content="black">
    <link rel="icon" href="data:;base64,=">
    <link rel="stylesheet" type="text/css" id="baseStyle" href="stage/build/export/base.css?{{ .Version }}"/>
    <link rel="stylesheet" type="text/css" id="themeDefaultStyle" href="appearance/themes/{{ .Theme }}/theme.css?{{ .Version }}"/>
    <link rel="stylesheet" type="text/css" id="themeStyle" href="appearance/themes/{{ .Theme }}/theme.css?{{ .Version }}"/>
    <title>{{ .Title }}{{ .TitleVersion }}</title>
    <style>
        body {font-family: var(--b3-font-family);background-color: var(--b3-theme-background);color: var(--b3-theme-on-background)}
        .b3-typography, .protyle-wysiwyg, .protyle-title {font-size:16px !important}
.b3-typography code:not(.hljs), .protyle-wysiwyg span[data-type~=code] { font-variant-ligatures: normal }
.li > .protyle-action {height:34px;line-height: 34px}
.protyle-wysiwyg [data-node-id].li > .protyle-action ~ .h1, .protyle-wysiwyg [data-node-id].li > .protyle-action ~ .h2, .protyle-wysiwyg [data-node-id].li > .protyle-action ~ .h3, .protyle-wysiwyg [data-node-id].li > .protyle-action ~ .h4, .protyle-wysiwyg [data-node-id].li > .protyle-action ~ .h5, .protyle-wysiwyg [data-node-id].li > .protyle-action ~ .h6 {line-height:34px;}
.protyle-wysiwyg [data-node-id].li > .protyle-action:after {height: 16px;width: 16px;margin:-8px 0 0 -8px}
.protyle-wysiwyg [data-node-id].li > .protyle-action svg {height: 14px}
.protyle-wysiwyg [data-node-id].li:before {height: calc(100% - 34px);top:34px}
.protyle-wysiwyg [data-node-id] [spellcheck] {min-height:26px;}
.protyle-wysiwyg [data-node-id] {}
.protyle-wysiwyg .li {min-height:34px}
.protyle-gutters button svg {height:26px}
        #helloPanel {
  border: 1px rgb(189, 119, 119) dashed;
}

.plugin-sample__custom-tab {
  background-color: var(--b3-theme-background);
  height: 100%;
  width: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
}
.plugin-sample__custom-dock {
  display: flex;
  justify-content: center;
  align-items: center;
}
.plugin-sample__time {
  background: var(--b3-card-info-background);
  border-radius: 4px;
  padding: 2px 8px;
}.config__panel.svelte-1fh5pey.svelte-1fh5pey.svelte-1fh5pey{height:100%}.config__panel.svelte-1fh5pey>ul.svelte-1fh5pey>li.svelte-1fh5pey{padding-left:1rem}

    </style>
</head>
<body>
<div class="protyle-wysiwyg protyle-wysiwyg--attr"
style="max-width: 800px;margin: 0 auto;"
id="preview">
{{ .Content }}
<script src="appearance/icons/material/icon.js?{{ .Version }}"></script>
<script src="stage/build/export/protyle-method.js?{{ .Version }}"></script>
<script src="stage/protyle/js/lute/lute.min.js?{{ .Version }}"></script>
<script>
    window.siyuan = {
      config: {
        appearance: { mode: 0, codeBlockThemeDark: "base16/dracula", codeBlockThemeLight: "github" },
        editor: {
          codeLineWrap: true,
          fontSize: 16,
          codeLigatures: true,
          plantUMLServePath: "https://www.plantuml.com/plantuml/svg/~1",
          codeSyntaxHighlightLineNum: true,
          katexMacros: JSON.stringify({}),
        }
      },
      languages: {copy:"复制"}
    };
    const previewElement = document.getElementById('preview');
    Protyle.highlightRender(previewElement, "stage/protyle");
    Protyle.mathRender(previewElement, "stage/protyle", false);
    Protyle.mermaidRender(previewElement, "stage/protyle");
    Protyle.flowchartRender(previewElement, "stage/protyle");
    Protyle.graphvizRender(previewElement, "stage/protyle");
    Protyle.chartRender(previewElement, "stage/protyle");
    Protyle.mindmapRender(previewElement, "stage/protyle");
    Protyle.abcRender(previewElement, "stage/protyle");
    Protyle.htmlRender(previewElement);
    Protyle.plantumlRender(previewElement, "stage/protyle");
    document.querySelectorAll(".protyle-action__copy").forEach((item) => {
      item.addEventListener("click", (event) => {
            let text = item.parentElement.nextElementSibling.textContent.trimEnd();
            text = text.replace(/ /g, " "); // Replace non-breaking spaces with normal spaces when copying
            navigator.clipboard.writeText(text);
            event.preventDefault();
            event.stopPropagation();
      })
    });
</script></body></html>
`
