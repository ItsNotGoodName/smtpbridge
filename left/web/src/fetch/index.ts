import { Ref, ref, unref, isRef, watchEffect } from 'vue'
import { IRequest, IResponse } from "../api"

export function useFetch<T>(req: Ref<IRequest<T>> | IRequest<T>, config = {} as { skip?: boolean, watch?: boolean }) {
  const data = ref<T | null>(null)
  const error = ref<string | null>(null)
  const loading = ref(true)

  function doFetch() {
    data.value = null
    error.value = null
    loading.value = true
    let r = unref(req)
    fetch(r.url, {
      method: r.method,
    }).then(res => res.json())
      .then((json: IResponse<T>) => {
        if (json.ok) {
          data.value = json.data as any
        } else {
          error.value = json.error!.message
        }
      })
      .catch((err) => (error.value = err.message))
      .finally(() => loading.value = false)
  }

  if (config.watch && isRef(req)) {
    watchEffect(doFetch)
  } else if (!config.skip) {
    doFetch()
  }

  return { data, error, loading, fetch: doFetch }
}
