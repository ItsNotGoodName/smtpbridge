// Fork of https://htmx.org/extensions/loading-states/ with the following changes
// - The loading container does not bind to document.body
// - Added data-loading-hide-class attribute which adds class to the element when a request is made
// - Import htmx so @rollup/plugin-inject is not needed

import htmx from "htmx.org"

let loadingStatesUndoQueue = []

function loadingStateContainer(target) {
  return htmx.closest(target, '[data-loading-states]') || target
}

function mayProcessUndoCallback(target, callback) {
  if (document.body.contains(target)) {
    callback()
  }
}

function mayProcessLoadingStateByPath(elt, requestPath) {
  const pathElt = htmx.closest(elt, '[data-loading-path]')
  if (!pathElt) {
    return true
  }

  return pathElt.getAttribute('data-loading-path') === requestPath
}

function queueLoadingState(sourceElt, targetElt, doCallback, undoCallback) {
  const delayElt = htmx.closest(sourceElt, '[data-loading-delay]')
  if (delayElt) {
    const delayInMilliseconds =
      delayElt.getAttribute('data-loading-delay') || 200
    const timeout = setTimeout(() => {
      doCallback()

      loadingStatesUndoQueue.push(() => {
        mayProcessUndoCallback(targetElt, () => undoCallback())
      })
    }, delayInMilliseconds)

    loadingStatesUndoQueue.push(() => {
      mayProcessUndoCallback(targetElt, () => clearTimeout(timeout))
    })
  } else {
    doCallback()
    loadingStatesUndoQueue.push(() => {
      mayProcessUndoCallback(targetElt, () => undoCallback())
    })
  }
}

function getLoadingStateElts(loadingScope, type, path) {
  return Array.from(htmx.findAll(loadingScope, `[${type}]`)).filter(
    (elt) => mayProcessLoadingStateByPath(elt, path)
  )
}

function getLoadingTarget(elt) {
  if (elt.getAttribute('data-loading-target')) {
    return Array.from(
      htmx.findAll(elt.getAttribute('data-loading-target'))
    )
  }
  return [elt]
}

htmx.defineExtension('loading-states', {
  onEvent: function(name, evt) {
    if (name === 'htmx:beforeRequest') {
      const container = loadingStateContainer(evt.target)

      const loadingStateTypes = [
        'data-loading',
        'data-loading-class',
        'data-loading-class-remove',
        'data-loading-disable',
        'data-loading-aria-busy',
        'data-loading-hide-class',
      ]

      let loadingStateEltsByType = {}

      loadingStateTypes.forEach((type) => {
        loadingStateEltsByType[type] = getLoadingStateElts(
          container,
          type,
          evt.detail.pathInfo.requestPath
        )
      })

      loadingStateEltsByType['data-loading'].forEach((sourceElt) => {
        getLoadingTarget(sourceElt).forEach((targetElt) => {
          queueLoadingState(
            sourceElt,
            targetElt,
            () =>
            (targetElt.style.display =
              sourceElt.getAttribute('data-loading') ||
              'inline-block'),
            () => (targetElt.style.display = 'none')
          )
        })
      })

      loadingStateEltsByType['data-loading-class'].forEach(
        (sourceElt) => {
          const classNames = sourceElt
            .getAttribute('data-loading-class')
            .split(' ')

          getLoadingTarget(sourceElt).forEach((targetElt) => {
            queueLoadingState(
              sourceElt,
              targetElt,
              () =>
                classNames.forEach((className) =>
                  targetElt.classList.add(className)
                ),
              () =>
                classNames.forEach((className) =>
                  targetElt.classList.remove(className)
                )
            )
          })
        }
      )

      loadingStateEltsByType['data-loading-class-remove'].forEach(
        (sourceElt) => {
          const classNames = sourceElt
            .getAttribute('data-loading-class-remove')
            .split(' ')

          getLoadingTarget(sourceElt).forEach((targetElt) => {
            queueLoadingState(
              sourceElt,
              targetElt,
              () =>
                classNames.forEach((className) =>
                  targetElt.classList.remove(className)
                ),
              () =>
                classNames.forEach((className) =>
                  targetElt.classList.add(className)
                )
            )
          })
        }
      )

      loadingStateEltsByType['data-loading-disable'].forEach(
        (sourceElt) => {
          getLoadingTarget(sourceElt).forEach((targetElt) => {
            queueLoadingState(
              sourceElt,
              targetElt,
              () => (targetElt.disabled = true),
              () => (targetElt.disabled = false)
            )
          })
        }
      )

      loadingStateEltsByType['data-loading-aria-busy'].forEach(
        (sourceElt) => {
          getLoadingTarget(sourceElt).forEach((targetElt) => {
            queueLoadingState(
              sourceElt,
              targetElt,
              () => (targetElt.setAttribute("aria-busy", "true")),
              () => (targetElt.removeAttribute("aria-busy"))
            )
          })
        }
      )

      loadingStateEltsByType['data-loading-hide-class'].forEach(
        (sourceElt) => {
          const classNames = sourceElt
            .getAttribute('data-loading-hide-class')
            .split(' ')

          getLoadingTarget(sourceElt).forEach((targetElt) => {
            queueLoadingState(
              sourceElt,
              targetElt,
              () =>
                classNames.forEach((className) => {
                  targetElt.classList.add(className)
                  targetElt.classList.add(className)
                }),
              () => { },
            )
          })
        }
      );
    }

    if (name === 'htmx:beforeOnLoad') {
      while (loadingStatesUndoQueue.length > 0) {
        loadingStatesUndoQueue.shift()()
      }
    }
  },
})

