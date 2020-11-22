import Vue from 'vue'
import { fromEvent, interval } from 'rxjs'
import { throttle } from 'rxjs/operators'

Vue.directive('resized', {
  inserted(el, binding) {
    const subscription = fromEvent(window, 'resize')
      // eslint-disable-next-line @typescript-eslint/no-unused-vars
      .pipe(throttle(() => interval(500)))
      .subscribe((evt) => {
        if (binding.value(evt, el)) {
          subscription?.unsubscribe()
        }
      })
  },
})
