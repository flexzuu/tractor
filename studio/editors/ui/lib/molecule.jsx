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

export function DropdownMenu(initial) {
    let opened = initial.attrs.opened || false;

    function toggle(e) {
        if (opened) {
            opened = false;
            return;
        }

        opened = true
        e.stopPropagation();
        document.body.addEventListener("click", clickoff);
    }

    function clickoff() {
        opened = false;
        m.redraw();
        document.body.removeEventListener("click", clickoff);
    }

    return {
        view: function (vnode) {
            return <div class="DropdownMenu" onclick={toggle}>
                {vnode.children}
                {opened && <Menu class={vnode.attrs.class} items={vnode.attrs.items} />}
            </div>;
        }
    }
}

export function Menu(initial) {
    return {
        view: function (vnode) {
            let items = vnode.attrs.items || [];
            return <ul class={"Menu right-0 absolute py-1 text-white w-32 rounded " + vnode.attrs.class} style={{ background: "#383838" }}>
                {items.map((item) =>
                    (item && item.label) ?
                        <li class="pl-5" onclick={item.onclick}>{item.label}</li> :
                        <li class="pl-5 my-1" style={{ "border-bottom": "3px solid #616161" }}></li>
                )}
            </ul>;
        }
    }
}