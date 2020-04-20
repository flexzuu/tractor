import * as atom from '/views/ui/lib/atom.js';

export function TextInput(initial) {
    return {
        view: function (vnode) {
            return <InputBox><input
                type="text"
                class="flex-auto"
                onchange={vnode.attrs.onchange}
                style={{ width: "0px", maxWidth: "100%", minWidth: "20%" }}
                value={vnode.attrs.value} />
            </InputBox>;
        }
    }
}

export function PasswordInput(initial) {
    return {
        view: function (vnode) {
            // TODO: copy textinput?
            return <InputBox>
                <input type="password"
                    class="flex-auto"
                    onchange={vnode.attrs.onchange}
                    style={{ width: "0px", maxWidth: "100%", minWidth: "20%" }}
                    autocomplete="password"
                    value={vnode.attrs.value}
                />
            </InputBox>;
        }
    }
}

export function NumberInput(initial) {
    return {
        view: function (vnode) {
            return <InputBox>
                <input type="number"
                    class="flex-auto"
                    onchange={vnode.attrs.onchange}
                    value={vnode.attrs.value}
                    style={{ width: "0px", maxWidth: "100%", minWidth: "20%" }}
                />
            </InputBox>;
        }
    }
}

export function SliderInput(initial) {
    let value = initial.attrs.value || 0;
    let min = initial.attrs.min || 0;
    let max = initial.attrs.max || 100;
    let step = initial.attrs.step || 1;
    return {
        view: function (vnode) {
            return <InputBox transparent><atom.Slider onchange={vnode.attrs.onchange} min={min} max={max} value={value} step={step} /></InputBox>;
        }
    }
}

export function KnobInput(initial) {
    let value = initial.attrs.value || 0;
    let min = initial.attrs.min || 0;
    let max = initial.attrs.max || 100;
    let step = initial.attrs.step || 1;
    let sensitivity = initial.attrs.sensitivity || 1;
    return {
        view: function (vnode) {
            return <InputBox transparent>
                <atom.Knob
                    value={value}
                    min={min}
                    max={max}
                    step={step}
                    sensivity={sensitivity}
                    onchange={vnode.attrs.onchange} />
            </InputBox>;
        }
    }
}

export function SelectInput(initial) {
    let value = initial.attrs.value;
    return {
        view: function (vnode) {
            let children = vnode.children.slice();
            children.forEach((el) => {
                let v = el.text;
                if (el.attrs && el.attrs.value) {
                    v = el.attrs.value;
                }
                if (v == value) {
                    if (!el.attrs) {
                        el.attrs = {};
                    }
                    el.attrs.selected = true;
                }
            })
            return <InputBox>
                <select class="w-full" onchange={vnode.attrs.onchange}>
                    {children}
                </select>
            </InputBox>
        }
    }
}

export function ColorInput(initial) {
    let color = initial.attrs.value || "#ffffff";
    let wellStyle = {
        width: "2.25rem",
        borderRadius: "4px 0 0 4px",
        height: "1.75rem",
        marginLeft: "-0.5rem",
        marginRight: "0.5rem"
    };
    return {
        view: function (vnode) {
            function select(e) {
                vnode.dom.querySelector("input[type=color]").click();
            }
            function updateFromPicker(e) {
                color = vnode.dom.querySelector("input[type=color]").value;
                if (vnode.attrs.onchange) {
                    vnode.attrs.onchange(e);
                }
            }
            function updateFromTextbox(e) {
                color = vnode.dom.querySelector("input[type=text]").value;
                if (vnode.attrs.onchange) {
                    vnode.attrs.onchange(e);
                }
            }
            let style = Object.assign({}, wellStyle);
            style["backgroundColor"] = color;
            return (
                <InputBox>
                    <div style={style} onclick={select}>&nbsp;</div>
                    <input class="w-full" type="text" oninput={updateFromTextbox} value={color} />
                    <input class="hidden" type="color" onchange={updateFromPicker} value={color} />
                    <atom.Icon onclick={select} fa="fas fa-eye-dropper"></atom.Icon>
                </InputBox>
            );
        }
    }
}

export function ReferenceInput(initial) {
    let hover = false;
    return {
        view: function (vnode) {
            function mouseover(e) {
                hover = true;
            }
            function mouseout(e) {
                hover = false;
            }
            return (
                <InputBox onmouseover={mouseover} onmouseout={mouseout}>
                    <i class={"w-7 mr-2 " + vnode.attrs.icon}></i>
                    <input
                        class="w-full"
                        type="text"
                        onchange={vnode.attrs.onchange}
                        value={vnode.attrs.value}
                        placeholder={vnode.attrs.placeholder}
                    />
                    {!hover && <atom.Icon fa="fas fa-asterisk"></atom.Icon>}
                    {hover && <atom.Icon fa="fas fa-times-circle"></atom.Icon>}
                </InputBox>
            );
        }
    }
}

export function TimeInput(initial) {
    return {
        view: function (vnode) {
            return (
                <InputBox>
                    <input type="time" onchange={vnode.attrs.onchange} value={vnode.attrs.value} required />
                    <atom.Icon fa="far fa-clock"></atom.Icon>
                </InputBox>
            );
        }
    }
}

export function DateInput(initial) {
    return {
        view: function (vnode) {
            return (
                <InputBox>
                    <input type="date"
                        class="flex-auto"
                        style={{ width: "0px", maxWidth: "100%", minWidth: "20%", marginRight: "-1rem" }}
                        required
                        onchange={vnode.attrs.onchange}
                        value={vnode.attrs.value}
                    />
                    <atom.Icon class="pointer-events-none" fa="fas fa-calendar-day"></atom.Icon>
                </InputBox>
            );
        }
    }
}

export function CheckboxInput(initial) {
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
                onclick: onclick,
                label: vnode.attrs.label,
                onchange: vnode.attrs.onchange
            };
            if (checked) {
                attrs["checked"] = "checked";
            }
            return <InputBox transparent>{m(atom.Checkbox, attrs)}</InputBox>;
        }
    }
}


/////////////////////



export function InputBox(initial) {
    return {
        view: function (vnode) {
            let style = {}
            if (vnode.attrs.transparent) {
                style["background"] = "transparent";
            }
            return <form
                class="form-input"
                onmouseout={vnode.attrs.onmouseout}
                onmouseover={vnode.attrs.onmouseover}
                style={style}>
                {vnode.children}
            </form>;
        }
    }
}