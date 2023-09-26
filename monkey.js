// ==UserScript==
// @name         Train Tools
// @namespace    http://tampermonkey.net/
// @version      0.1
// @description  try to take over the world!
// @author       You
// @match        https://kyfw.12306.cn/*/*
// @icon         data:image/gif;base64,R0lGODlhAQABAAAAACH5BAEKAAEALAAAAAABAAEAAAICTAEAOw==
// @grant        none
// ==/UserScript==

(function() {
    'use strict';
// 创建一个 iframe 元素
    var iframe = document.createElement('iframe');

// 设置 iframe 样式
    iframe.style.position = 'fixed'; // 或 'absolute'，具体根据需要选择
    iframe.style.top = '0px';       // 距离顶部的距离
    iframe.style.left = '0px';      // 距离左侧的距离
    iframe.style.width = '100%';     // 宽度
    iframe.style.height = '100%';    // 高度
    iframe.style.zIndex = "9999";

// 其他样式设置，如背景颜色、边框等
    iframe.style.backgroundColor = 'white';
    iframe.style.border = '1px solid #ccc';

// 设置 iframe 内容，这部分代码保持不变
    iframe.srcdoc = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8">
    <link rel="icon" href="/favicon.ico">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Travel</title>
  </head>
  <body>
    <div id="app"></div>
    <script type="module" src="/src/main.ts"></script>
  </body>
</html>
    
    `;

// 将 iframe 添加到当前页面
    document.body.appendChild(iframe);
})();