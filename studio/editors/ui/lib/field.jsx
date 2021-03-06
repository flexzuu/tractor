import * as atom from '/views/ui/lib/atom.js';
import * as form from '/views/ui/lib/form.js';
import * as molecule from '/views/ui/lib/molecule.js';

function getPathTo(element) {
    if (!element) return "";
    if (element.id !== '')
        return "*[@id=" + element.id + "]";

    if (element === document.body)
        return element.tagName.toLowerCase();

    var ix = 0;
    var siblings = element.parentNode.childNodes;
    for (var i = 0; i < siblings.length; i++) {
        var sibling = siblings[i];

        if (sibling === element) return getPathTo(element.parentNode) + '/' + element.tagName.toLowerCase() + '[' + (ix + 1) + ']';

        if (sibling.nodeType === 1 && sibling.tagName === element.tagName) {
            ix++;
        }
    }
}

function dirname(path) {
    return path.match(/.*\//) || "";
}

// required for useLocalStorage
function uniqueIdent(vnode, suffix) {
    vnode.state.ident = [dirname(getPathTo(vnode.dom)) + vnode.tag.name, vnode.key, suffix].join(" ").replace(/ /g, "-");
    m.redraw();
}

// requires uniqueIdent to called in oncreate
function useLocalStorage(vnode, initial) {
    const id = vnode.state.ident;
    let value = JSON.parse(localStorage.getItem(`--useLocalStorage-${id}`)) || initial;
    let setter = function (v) {
        localStorage.setItem(`--useLocalStorage-${id}`, JSON.stringify(v));
    }
    return [value, setter];
}

export function ComponentPanel(initial) {
    let label = initial.attrs.label;
    let expanded = initial.attrs.expanded;
    return {
        oncreate: function (vnode) {
            uniqueIdent(vnode, label);
        },
        view: function (vnode) {
            var setExpanded;
            [expanded, setExpanded] = useLocalStorage(vnode, vnode.attrs.expanded);
            function toggle() {
                if (expanded) {
                    expanded = false;
                } else {
                    expanded = true;
                }
                setExpanded(expanded);
            }
            let bottomPadding = "pb-1";
            if (expanded) {
                bottomPadding = "pb-4"
            }
            let expanderMargin = "mb-2";
            if (expanded) {
                expanderMargin = "mb-4"
            }
            return (
                <div class={"ComponentPanel flex flex-col w-full my-1 pl-2 " + bottomPadding} style={{ borderBottom: "1px solid #404040" }}>
                    <molecule.Expander class={expanderMargin} expanded={expanded} onclick={toggle}>
                        <atom.Checkbox class="w-7 mr-2" />
                        <div onclick={toggle} title={label} class="label flex-grow h-6 truncate">{label}</div>
                        <molecule.DropdownMenu class="mr-2" items={vnode.attrs.menu}>
                            <atom.Icon class="mr-2 w-3" fa="fas fa-ellipsis-v" />
                        </molecule.DropdownMenu>
                    </molecule.Expander>
                    {/* {expanded && vnode.children.map((el) => <div class="my-1 mx-4">{el}</div>)} */}
                    {expanded && vnode.children}
                </div>
            );
        }
    }
}

export function CollectionItem(initial) {
    return {
        view: function (vnode) {
            let item = vnode.children;
            let idx = vnode.attrs.idx;
            function onclick(e) {
                vnode.attrs.onremove(idx);
            }
            if (vnode.attrs.removable) {
                item = <atom.Removable onclick={onclick}>{item}</atom.Removable>;
            }
            if (vnode.attrs.draggable) {
                item = <atom.Grippable>{item}</atom.Grippable>;
            }
            return <div class={"my-1 ml-4 " + vnode.attrs.class}>{item}</div>;
        }
    }
}

export function Nested(initial) {
    let label = initial.attrs.label;
    let expanded = initial.attrs.expanded;
    return {
        oncreate: function (vnode) {
            uniqueIdent(vnode, label);
        },
        view: function (vnode) {
            var setExpanded;
            [expanded, setExpanded] = useLocalStorage(vnode, vnode.attrs.expanded);
            function toggleExpander() {
                if (expanded) {
                    expanded = false;
                } else {
                    expanded = true;
                }
                setExpanded(expanded);
            }
            let expanderMargin = "mb-0";
            if (expanded) {
                expanderMargin = "mb-2"
            }
            return (
                <div class="Nested flex flex-col select-none pb-2">
                    <molecule.Expander expanded={expanded} class={expanderMargin} onclick={toggleExpander}>
                        <div onclick={toggleExpander} class="label flex-grow h-4">{label}</div>
                    </molecule.Expander>
                    {expanded && vnode.children.map((el) => <div class="my-1 ml-4">{el}</div>)}
                </div>
            );
        }
    }
}

export function Collection(initial) {
    let subtype = initial.attrs.subtype;
    let label = initial.attrs.label;
    let expanded = initial.attrs.expanded;
    let onadd = initial.attrs.onadd;
    let onremove = initial.attrs.onremove;
    let adding = false;
    var newValue;
    return {
        oncreate: function (vnode) {
            uniqueIdent(vnode, label);
        },
        view: function (vnode) {
            var setExpanded;
            [expanded, setExpanded] = useLocalStorage(vnode, expanded);
            function toggleExpander() {
                if (expanded) {
                    expanded = false;
                    adding = false;
                } else {
                    expanded = true;
                }
                setExpanded(expanded);
            }
            function toggleAdd(e) {
                if (adding) {
                    adding = false;
                } else {
                    adding = true;
                }
            }
            let expanderMargin = "mb-0";
            if (expanded) {
                expanderMargin = "mb-2"
            }
            function changeNewValue(e) {
                newValue = e.target.checked || e.target.value;
                if (e.target.type === "number") {
                    newValue = e.target.valueAsNumber;
                }
            }
            function add(e) {
                if (onadd) {
                    onadd(newValue);
                }
                newValue = undefined;
                adding = false;
            }
            return (
                <div class="Collection flex flex-col select-none">
                    <molecule.Expander expanded={expanded} class={expanderMargin} onclick={toggleExpander}>
                        <div onclick={toggleExpander} class="label flex-grow h-4">{label}</div>
                        <span class="mr-2 mt-1 text-xs">{vnode.children.length} items</span>
                        <atom.Icon class="mr-2" fa="fas fa-plus-circle" onclick={toggleAdd} />
                    </molecule.Expander>
                    {adding && <CollectionItem class="flex flex-col my-4">
                        <Input field={subtype} onchange={changeNewValue} value={newValue} />
                        <atom.Button class="mt-2" label="Add" onclick={add} />
                    </CollectionItem>}
                    {expanded && vnode.children.map((el, idx) => <CollectionItem idx={idx} onremove={onremove} removable draggable>{el}</CollectionItem>)}
                </div>
            );
        }
    }
}

export function Row(initial) {
    return {
        view: function (vnode) {
            const children = vnode.children.slice(0);
            const label = children.shift();
            return <div class="flex">
                <div class="mr-2" style={{ minWidth: "35%" }}>{label}</div>
                <div class="flex-grow flex">{children}</div>
            </div>;
        }
    }
}

export function LabeledField(initial) {
    return {
        view: function (vnode) {
            return <Row>
                <span class="text-sm">{vnode.attrs.label}</span>
                {vnode.children}
            </Row>;
        }
    }
}

export function KeyedField(initial) {
    return {
        view: function (vnode) {
            return <Row>
                <form.TextInput />
                {vnode.children}
            </Row>;
        }
    }
}

export function Input(initial) {
    return {
        view: function (vnode) {
            let field = vnode.attrs.field;
            let send = console.log;
            if (window.remoteCall) {
                send = window.remoteCall;
            }
            let onchange = vnode.attrs.onchange || function (e) {
                switch (e.target.type) {
                    case "checkbox":
                        send("setValue", { "Type": field.type, "Path": field.path, "Value": e.target.checked });
                        break;
                    case "number":
                        send("setValue", { "Type": field.type, "Path": field.path, "IntValue": e.target.valueAsNumber });
                        break;
                    default:
                        send("setValue", { "Type": field.type, "Path": field.path, "Value": e.target.value });
                }
            }
            switch (field.type) {
                case "string":
                    if (field.enum) {
                        return <form.SelectInput onchange={onchange} value={field.value}>
                            {field.enum.map((opt) => <option>{opt}</option>)}
                        </form.SelectInput>;
                    }
                    return <form.TextInput onchange={onchange} value={field.value} />
                case "boolean":
                    return <form.CheckboxInput onchange={onchange} checked={field.value} />
                case "number":
                    return <form.NumberInput onchange={onchange} value={field.value} />
                case "time":
                    return <form.TimeInput onchange={onchange} value={field.value} />
                case "date":
                    return <form.DateInput onchange={onchange} value={field.value} />
                case "color":
                    return <form.ColorInput onchange={onchange} value={field.value} />
                default:
                    if (field.type.startsWith("reference:")) {
                        var refType = field.type.split(":")[1];
                        function onset(path) {
                            send("setValue", { "Type": field.type, "Path": field.path, "RefValue": `${path}/${refType}` });
                        }
                        function onunset(path) {
                            send("setValue", { "Type": field.type, "Path": field.path, "Value": null });
                        }
                        return <form.ReferenceInput value={field.value} placeholder={refType} onset={onset} onunset={onunset} />;
                    } else {
                        return `Unknown field type: ${field.type}`;
                    }
            }
        }
    }
}

export function ComponentField(initial) {
    return {
        view: function (vnode) {
            let field = vnode.attrs.field || {};
            let subfields = field.fields || [];
            let nolabel = vnode.attrs.nolabel;

            let send = console.log;
            if (window.remoteCall) {
                send = window.remoteCall;
            }
            function collectionAdd(value) {
                send("addValue", { "Type": field.subtype.type, "Path": field.path, "Value": value });
            }
            function collectionRemove(value) {
                send("removeValue", { "Type": field.type, "Path": field.path, "IntValue": value });
            }
            switch (field.type) {
                case "struct":
                    return (
                        <Nested key={vnode.key} label={field.name}>
                            {subfields.map((f) => <ComponentField key={f.name} field={f} />)}
                        </Nested>);
                case "map":
                    return (
                        <Nested key={vnode.key} label={field.name}>
                            {subfields.map((f) => <KeyedField key={f.name} name={f.name}><Input field={f} /></KeyedField>)}
                        </Nested>);
                case "array":
                    return (
                        <Collection key={vnode.key} label={field.name} subtype={field.subtype} onadd={collectionAdd} onremove={collectionRemove}>
                            {subfields.map((f) => <ComponentField key={f.name} field={f} nolabel={true} />)}
                        </Collection>);
                default:
                    if (nolabel === true) {
                        return <Input key={vnode.key} field={field} />
                    }
                    return <LabeledField key={vnode.key} label={field.name}><Input field={field} /></LabeledField>
            }
        }
    }
}
