import "./index.css"

// ------------- htmx

import "htmx.org"

const csrfToken = (document.getElementsByName("gorilla.csrf.Token").item(0) as HTMLMetaElement).content

document.body.addEventListener('htmx:configRequest', function(evt: any) {
  evt.detail.headers['X-CSRF-Token'] = csrfToken;
});

// ------------- Toastify

import Toastify from 'toastify-js'

document.body.addEventListener('htmx:afterRequest', function(evt: any) {
  if (evt.detail.failed) {
    const content = document.createElement("div")
    content.textContent = evt.detail.xhr.responseText || evt.detail.xhr.statusText
    content.className = "flex-1"

    Toastify({
      text: evt.detail.xhr.responseText || evt.detail.xhr.statusText,
      node: content,
      duration: 3000,
      close: true,
      className: "alert alert-error flex flex-row",
      gravity: "bottom",
      position: "center",
      stopOnFocus: true,
    }).showToast();
  }
})

// ------------- Shoelace

// 24 Kb to format dates in local time
import '@shoelace-style/shoelace/dist/components/format-date/format-date.js';
