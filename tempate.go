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
    {{ .MiniMenuStyle }}
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
        <style>
    {{ .CustomCSS }}
    </style>
</head>
<body>
<div class="protyle-wysiwyg protyle-wysiwyg--attr"
style="max-width: {{ .PageWide }};margin: 0 auto;"
id="preview">
<div>
  <div id="toc-container">
  </div>
  <button id="toggle-button"></button>
</div>
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
</script>

{{ .MiniMenuScript }}
</body></html>
`
const access_key_html = `
<!DOCTYPE html>
<html>
<head>
  <title>访问码</title>
</head>
<body>

<script>
 function getAccessKey() {
    let accessKey = prompt("请输入访问码：\n Please enter access code:");
    if (!accessKey) {
      getAccessKey(); 
      return;
    } 
    window.location.href =  window.location.origin + window.location.pathname + '?access_key=' + encodeURIComponent(accessKey);
  }
  getAccessKey(); 
</script>

</body>
</html>
 
`
const mini_menu_script = `
  <script>
  const divs = document.querySelectorAll('div[data-subtype]');
  const headings = [];

  divs.forEach(div => {
  const subtype = div.getAttribute('data-subtype');
  const nodeId = div.getAttribute('data-node-id'); // 获取 data-node-id 属性
  if (['h1', 'h2', 'h3', 'h4', 'h5', 'h6'].includes(subtype)) {
    const title = div.firstElementChild.textContent;
    headings.push({
      nodeId: nodeId, // 使用 data-node-id
      subtype: subtype,
      title: title
    });
  }
});

let tocHtml = '<ul>';
headings.forEach(heading => {
  const indent = (parseInt(heading.subtype.replace('h', '')) - 1) * 20; // 计算缩进值
  tocHtml += ` + "`<li style=\"padding-left: ${indent}px\"><a href=\"#${heading.nodeId}\">${heading.title}</a></li>`; // 使用 data-node-id " + `
});
tocHtml += '</ul>';

const tocContainer = document.getElementById('toc-container');
tocContainer.innerHTML = tocHtml;

const tocLinks = document.querySelectorAll('#toc-container a');
tocLinks.forEach(link => {
  link.addEventListener('click', (event) => {
    event.preventDefault();
    const targetId = link.getAttribute('href').substring(1); // 去掉 '#' 符号
    const targetDiv = document.querySelector(` + "`[data-node-id=\"${targetId}\"]`" + `); // 使用 data-node-id 查询元素
    
        // 调试信息
    // console.log('Scrolled to:', targetId);
    if (targetDiv) {
      // console.log('Target element found:', targetDiv);
      targetDiv.scrollIntoView({ behavior: 'smooth' });
    } else {
      console.error('Target element not found for:', targetId);
    }

  });
});
 var tocVisible = true; // 菜单初始状态为可见

document.getElementById('toggle-button').onclick = function () {
  var tocContainer = document.getElementById('toc-container');
  var toggleButton = document.getElementById('toggle-button');
  if (tocVisible) {
    tocContainer.style.transform = 'translateY(-50%) translateX(-400%)';
    toggleButton.style.borderWidth = '10px 0 10px 10px'; // 修改箭头方向
    toggleButton.style.borderColor = 'transparent transparent transparent #007bff'; // 修改箭头颜色
  } else {
    tocContainer.style.transform = 'translateY(-50%) translateX(0)';
    toggleButton.style.borderWidth = '10px 10px 10px 0';
    toggleButton.style.borderColor = 'transparent #007bff transparent transparent'; // 还原箭头颜色
  }
  tocVisible = !tocVisible; // 切换菜单状态
}
</script>
`
const mini_menu_style = `
<style>
#toc-container {
  position: fixed;
  top: 50%;
  left: 20px;
  transform: translateY(-50%);
  padding: 10px;
  background-color: #fff;
  border-radius: 5px;
  z-index: 100; /* 确保菜单在其他元素上方 */
}

#toc-container ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

#toc-container li a {
  display: block;
  padding: 5px 10px;
  color: #333;
  text-decoration: none;
  transition: all 0.2s ease;
}

#toc-container li a:hover {
  color: #333;
}
#toggle-button {
  position: fixed;
  top: 50%;
  width: 0;
  height: 0;
  left:10px;
  border-style: solid;
  border-width: 10px 10px 10px 0;
  border-color: transparent #007bff transparent transparent;
  cursor: pointer;
  transform: translateY(-50%);
  transition: left 0.3s ease, transform 0.3s ease;
  background: transparent;
  z-index: 101;
}
</style>`
