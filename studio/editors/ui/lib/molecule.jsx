import * as atom from '/views/ui/lib/atom.js';

export function Expander(initial) {
    return {
        view: function (vnode) {
            let expanded = vnode.attrs.expanded || false;
            return <div class={"Expander flex h-4 " + vnode.attrs.class}>
                <atom.Icon onclick={vnode.attrs.onclick} class="w-4 text-center" fa={`fas fa-caret-${(expanded) ? 'down' : 'right'}`} />
                {vnode.children}
            </div>;
        }
    }
}