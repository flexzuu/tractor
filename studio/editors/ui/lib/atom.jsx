
export function Button(initial) {
    let state = "up";
    let style = {
        backgroundColor: "#404040"
    }
    return {
        view: function (vnode) {
            function press(e) {
                state = "down";
            }
            function release(e) {
                state = "up";
            }
            return <div class={"Button rounded " + state + " " + vnode.attrs.class} style={{ border: "2px solid black" }}>
                <button onclick={vnode.attrs.onclick} onmousedown={press} onmouseup={release} blur={release} onmouseout={release} class="w-full h-8 rounded px-2" style={style}>
                    {vnode.attrs.label || "Button"}
                </button>
            </div>;
        }
    }
}

export function Label(initial) {
    return {
        view: function (vnode) {
            return m("span", vnode.attrs, vnode.children);
        }
    }
}

export function Grip(initial) {
    return {
        view: function (vnode) {
            let attrs = vnode.attrs;
            let style = {
                backgroundColor: "transparent",
                backgroundImage: "radial-gradient(#404040 50%, transparent 50%)",
                backgroundSize: "4px 4px",
                backgroundRepeat: "repeat",
                backgroundPosition: "0px 4px",
                width: "8px",
                height: "1.5rem",
                flex: "0 0 auto",
                marginTop: "0.125rem",
                marginRight: "0.25rem"
            }
            attrs.style = Object.assign(style, attrs.style || {});
            return m("div", attrs, vnode.children);
        }
    }
}


export function Icon(initial) {
    return {
        view: function (vnode) {
            return <div class={vnode.attrs.class} onclick={vnode.attrs.onclick}><i class={vnode.attrs.fa}></i></div>;
        }
    }
}

export function Grippable(initial) {
    return {
        view: function (vnode) {
            return <div class="Grippable flex"><Grip /><div class="flex-grow">{vnode.children}</div></div>;
        }
    }
}

export function Removable(initial) {
    return {
        view: function (vnode) {
            return <div class="Removable">
                <Icon class="float-right mx-2" fa="fas fa-times-circle"></Icon>
                <div class="mr-8">{vnode.children}</div>
            </div>;
        }
    }
}

export function Indicator(initial) {
    return {
        view: function (vnode) {
            let active = vnode.attrs.active;
            let color = vnode.attrs.color || "#ffff00";
            let style = {
                outer: {
                    width: "1.5rem",
                    height: "1.5rem",
                    borderRadius: "50%",
                    backgroundColor: "#404040"
                },
                inner: {
                    position: "relative",
                    left: "0.38rem",
                    top: "0.38rem",
                    width: "0.75rem",
                    height: "0.75rem",
                    borderRadius: "50%",
                    backgroundColor: color,
                    opacity: (active) ? 1 : 0.25
                }
            }
            if (color !== "black" && active) {
                style.inner["filter"] = `drop-shadow(0 0 0.25rem ${color})`;
            }
            return <div class="Indicator" style={style.outer}><div style={style.inner}></div></div>;
        }
    }
}

export function Knob(initial) {
    let value = initial.attrs.value || 0;
    let min = initial.attrs.min || 0;
    let max = initial.attrs.max || 100;
    let step = initial.attrs.step || 1;
    let sensitivity = initial.attrs.sensitivity || 1;
    let color = "white";
    let style = {
        knob: {
            width: "1.5rem",
            height: "1.5rem",
            background: "#404040",
            borderRadius: "50%",
            position: "relative",
            top: "0.25rem"
        },
        indicator: {
            height: "0.8rem",
            width: "0.2rem",
            position: "relative",
            left: "50%",
            top: "50%",
            marginLeft: "-0.1rem",
            borderRadius: "2px",
            transform: "rotate(180deg)",
            transformOrigin: "50% 0%"
        }
    };
    return {
        view: function (vnode) {
            function onmousedown(e) {
                let lastX = e.pageX;
                let lastY = e.pageY;

                color = "yellow";

                function onmousemove(e) {
                    let offsetX = e.pageX - lastX;
                    let offsetY = e.pageY - lastY;
                    let diff = (offsetX - offsetY) * sensitivity;
                    value = Math.min(max, Math.max(min, value + diff));
                    m.redraw();
                    lastX = e.pageX;
                    lastY = e.pageY;
                }

                function onmouseup(e) {
                    color = "white";
                    m.redraw();

                    document.body.removeEventListener("mouseup", onmouseup);
                    document.body.removeEventListener("mousemove", onmousemove);
                    return false;
                }

                document.body.addEventListener("mouseup", onmouseup);
                document.body.addEventListener("mousemove", onmousemove);
                return false;
            }
            let indicator = Object.assign({}, style.indicator);
            indicator["backgroundColor"] = color;
            let rot = (value * 300 / (max - min)) + 30;
            indicator["transform"] = `rotate(${rot}deg)`;
            return (
                <div style={style.knob} onmousedown={onmousedown}>
                    <div style={indicator}></div>
                    <input type="range" onchange={vnode.attrs.onchange} class="hidden" min={min} max={max} value={value} step={step} />
                </div>
            );
        }
    }
}

export function Slider(initial) {
    let value = initial.attrs.value || 0;
    let min = initial.attrs.min || 0;
    let max = initial.attrs.max || 100;
    let step = initial.attrs.step || 1;
    return {
        // see also range.css
        view: function (vnode) {
            return <input type="range" onchange={vnode.attrs.onchange} min={min} max={max} value={value} step={step} />;
        }
    }
}

export function Checkbox(initial) {
    let checked = initial.attrs.checked || false;
    let onclick = initial.attrs.onclick || function (e) {
        if (checked) {
            checked = false;
        } else {
            checked = true;
        }
    };
    return {
        view: function (vnode) {
            let attrs = {
                type: "checkbox",
                style: { // see also form.css
                    "-webkit-appearance": "none",
                    appearance: "none",
                    backgroundColor: "#404040",
                    width: "1rem",
                    height: "1rem",
                    marginTop: "0.25rem",
                },
                onclick: onclick,
                onchange: vnode.attrs.onchange
            };
            if (checked) {
                attrs["checked"] = "checked";
            }
            return (
                <div class={"inline-flex " + vnode.attrs.class}>
                    {m("input", attrs)}<span></span>
                    {vnode.attrs.label && <div class="ml-2">{vnode.attrs.label}</div>}
                </div>
            );
        }
    }
}