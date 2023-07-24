import "htmx.org"
import "htmx.org/dist/ext/loading-states.js"
import "@picocss/pico";
import "./pico.css"
import 'uno.css'

function getCookieByKey(key: string): string {
  const cookies = document.cookie.split(';');

  for (let i = 0; i < cookies.length; i++) {
    const cookie = cookies[i].trim();
    if (cookie.startsWith(`${key}=`)) {
      return decodeURIComponent(cookie.substring(key.length + 1));
    }
  }

  return "";
}

document.body.addEventListener('htmx:configRequest', function(evt: any) {
  const csrf = getCookieByKey("csrf_")
  if (csrf) {
    evt.detail.headers['X-CSRF-Token'] = csrf;
  }
});
