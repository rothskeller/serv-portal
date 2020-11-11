// This method, called from the setup() function of a container component,
// makes that component 'provide' a property named 'containerSize'.  Descendants
// of the container can inject it to get the size of the nearest ancestor
// container whose setup() function called this method.  The property value is
// a reactive object with 'w' and 'h' elements, each measured in rem.  In
// browsers that do not support ResizeObserver, 'w' and 'h' will be zero.
import { provide, reactive, onMounted, getCurrentInstance, onBeforeUnmount } from "vue"

export type Size = {
    w: number
    h: number
}

const observing: Map<Element, Size> = new Map()
const rem = parseFloat(getComputedStyle(document.documentElement).fontSize)

let ro: ResizeObserver
try {
    ro = new ResizeObserver(entries => {
        for (const entry of entries) {
            console.log('ro entry', entry)
            const size = observing.get(entry.target)
            if (size) {
                let bb = entry.borderBoxSize
                if (Array.isArray(bb)) bb = bb[0]
                size.w = bb.inlineSize / rem
                size.h = bb.blockSize / rem
                console.log('set size to', size.w, size.h)
            }
        }
    })
} catch (e) {
    console.warn('ResizeObserver could not be created')
}

function observe(el: Element, size: Size) {
    if (!observing.has(el)) ro.observe(el)
    observing.set(el, size)
}
function unobserve(el: Element, size: Size) {
    if (observing.has(el)) ro.unobserve(el)
    observing.delete(el)
}

// reactive and contains 'w' and 'h' values, both numbers in pixels.
export default () => {
    const size = reactive({ w: 0, h: 0 })
    provide('containerSize', size)
    if (!ro) return size
    onMounted(() => {
        const instance = getCurrentInstance()
        if (!instance) throw (new Error('no current instance in onMounted'))
        if (!instance.vnode) throw (new Error('no vnode in onMounted'))
        if (!instance.vnode.el) throw (new Error('no vnode.el in onMounted'))
        observe(instance.vnode.el as Element, size)
        console.log('observing size of', instance)
    })
    onBeforeUnmount(() => {
        const instance = getCurrentInstance()
        if (!instance) throw (new Error('no current instance in onMounted'))
        if (!instance.vnode) throw (new Error('no vnode in onMounted'))
        if (!instance.vnode.el) throw (new Error('no vnode.el in onMounted'))
        unobserve(instance.vnode.el as Element, size)
        console.log('unobserving size of', instance)
    })
    return size
}

// Make Typescript happy with the use of ResizeObserver.
declare class ResizeObserver {
    constructor(callback: ResizeObserverCallback)
    disconnect: () => void
    observe: (target: Element, options?: ResizeObserverObserveOptions) => void
    unobserve: (target: Element) => void
}
interface ResizeObserverObserveOptions {
    box?: "content-box" | "border-box"
}
type ResizeObserverCallback = (
    entries: ResizeObserverEntry[],
    observer: ResizeObserver,
) => void
interface ResizeObserverEntry {
    // The spec says that it returns an array, but apparently some browsers
    // return a single object.  We'll handle both.
    readonly borderBoxSize: ResizeObserverEntryBoxSize | Array<ResizeObserverEntryBoxSize>
    readonly contentBoxSize: ResizeObserverEntryBoxSize | Array<ResizeObserverEntryBoxSize>
    readonly contentRect: DOMRectReadOnly
    readonly target: Element
}
interface ResizeObserverEntryBoxSize {
    blockSize: number
    inlineSize: number
}
