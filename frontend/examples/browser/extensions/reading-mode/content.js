// Author: HyperparameterSI
const old = document.getElementById("golet-reading-mode");
if (old) old.remove();
else {
  const style = document.createElement("style");
  style.id = "golet-reading-mode";
  style.textContent =
    "main{max-width:900px!important;margin-inline:auto!important;line-height:1.9!important}";
  document.head.append(style);
}
