import * as field from '/views/ui/lib/field.js';
import * as atom from '/views/ui/lib/atom.js';
import * as molecule from '/views/ui/lib/molecule.js';
import * as form from '/views/ui/lib/form.js';
import * as client from '/views/ui/lib/client.js';

export function App(initial) {
    App.env = JSON.parse(atob(window.location.hash.slice(1)) || "{}");

    return {
        view: function (vnode) {
            return <Inspector />;
        }
    }
}

export function Inspector(initial) {
    let remote = { components: [] };
    var lastSelected;
    var node;
    let path = App.env.workspace;
    let session = new client.Session(path, (client) => {
        client.call("subscribe");
        window.remoteCall = (action, params) => {
            switch (action) {
                case "setValue":
                case "setExpression":
                case "callMethod":
                case "refreshObject":
                case "removeComponent":
                case "appendComponent":
                case "updateNode":
                case "reloadComponent":
                case "addDelegate":
                case "addValue":
                case "removeValue":
                    //console.log(action, params);
                    return client.call(action, params);
                case "edit":
                    window.parent.postMessage({ event: 'edit', path: params.path });
                    return;
                default:
                    throw "unknown action: " + action;
            }
        }
    });
    session.api.handle("shutdown", {
        "serveRPC": async (r, c) => {
            console.log("DEBUG: reload/shutdown received, reconnecting in 3s...");
            setTimeout(() => {
                console.log("reconnecting...")
                session.reconnect();
            }, 3000);
            r.return();
        }
    });
    session.api.handle("state", {
        "serveRPC": async (r, c) => {
            remote = await c.decode();
            console.log(remote);
            node = remote.nodes[remote.selectedNode || lastSelected];
            if (remote.selectedNode) {
                lastSelected = remote.selectedNode;
            }
            m.redraw();
            r.return();
        }
    });



    return {
        view: function (vnode) {
            if (node) {
                function buildMenu(component) {
                    let send = console.log;
                    if (window.remoteCall) {
                        send = window.remoteCall;
                    }
                    return [
                        {
                            label: "Reload", onclick: () => {
                                send("reloadComponent",
                                    { ID: node.id, Component: component.name })
                            }
                        },
                        {
                            label: "Edit", onclick: () => {
                                send("edit", { path: component.filepath })
                            }
                        },
                        {
                            label: "Remove", onclick: () => {
                                send("removeComponent",
                                    { ID: node.id, Component: component.name })
                            }
                        },
                    ];
                }

                return <section class="">
                    <ObjectHeader node={node} />
                    {node.components.map((c) => {
                        return m(field.ComponentPanel, { label: c.name, menu: buildMenu(c) }, (c.customUI) ?
                            m(CustomUI, { spec: c.customUI, fields: c.fields }) :
                            m(DefaultUI, { fields: c.fields })
                        );
                    })}
                </section>;
            }
            return <section>no node selected</section>;
        }
    }
}

export function DefaultUI(initial) {
    return {
        view: function (vnode) {
            let fields = vnode.attrs.fields || [];
            return fields.map((f, idx) =>
                <div class="my-1 mx-4">{m(field.ComponentField, { key: idx, field: f })}</div>);
        }
    }
}

export function CustomUI(initial) {
    let allowedElements = {
        "form.TextInput": form.TextInput,
        "form.PasswordInput": form.PasswordInput,
        "form.NumberInput": form.NumberInput,
        "form.SliderInput": form.SliderInput,
        "form.KnobInput": form.KnobInput,
        "form.SelectInput": form.SelectInput,
        "form.ColorInput": form.ColorInput,
        "form.ReferenceInput": form.ReferenceInput,
        "form.TimeInput": form.TimeInput,
        "form.DateInput": form.DateInput,
        "form.CheckboxInput": form.CheckboxInput,
        "atom.Button": atom.Button,
        "atom.Label": atom.Label,
        "atom.Grip": atom.Grip,
        "atom.Icon": atom.Icon,
        "atom.Grippable": atom.Grippable,
        "atom.Removable": atom.Removable,
        "atom.Indicator": atom.Indicator,
        "atom.Knob": atom.Knob,
        "atom.Slider": atom.Slider,
        "atom.Checkbox": atom.Checkbox,
        "field.ComponentField": field.ComponentField
    };
    function inflate(el, fields) {
        let tag = allowedElements[el.Name];
        if (!tag) {
            tag = el.Name;
        }
        el.Attrs = el.Attrs || {};
        if (el.Attrs["data-field"]) {
            el.Attrs["field"] = findField(fields, el.Attrs["data-field"]);
            delete el.Attrs["data-field"]
        }
        return m(tag, el.Attrs, (el.Children || []).map((c) => inflate(c, fields)))
    }
    return {
        view: function (vnode) {
            return vnode.attrs.spec.map((spec) => <div class="my-1 mx-4">{inflate(spec, vnode.attrs.fields)}</div>);
        }
    }
}

// TODO: this searches for a fields basename, not a relative path
function findField(fields, name) {
    for (const f of Object.keys(fields)) {
        const field = fields[f];
        if (field.name === name) {
            return field;
        }
        if (field.fields) {
            let sub = findField(field.fields, name);
            if (sub) {
                return sub;
            }
        }
    }
}

export function ObjectHeader(initial) {
    return {
        view: function (vnode) {
            var node;
            if (!vnode.attrs.node) {
                node = { id: "", name: "", path: "" };
            } else {
                node = vnode.attrs.node;
            }
            let parents = node.path.split("/").slice(1, -1);
            let send = console.log;
            if (window.remoteCall) {
                send = window.remoteCall;
            }
            function onchange(e) {
                send("updateNode", { "ID": node.id, "Name": e.target.value });
            }
            let menu = [{
                label: "Reload", onclick: () => {
                    send("refreshObject", { ID: node.id })
                }
            }]
            return <div class="flex w-full pl-3 pr-2 py-2" style={{ "borderBottom": "2px solid #404040" }}>
                <div class="self-end"><atom.Icon fa="fab fa-dev fa-3x" /></div>
                <div class="ml-2 flex-grow self-end" style={{ "maxWidth": "80%" }}>
                    <div class="text-xs flex flex-no-wrap breadcrumbs">
                        {parents.map((name) => <span>{name}</span>)}
                    </div>
                    <div class="w-full"><form.TextInput onchange={onchange} title={node.id} value={node.name} /></div>
                </div>
                <div class="ml-2 w-4 self-end">
                    <molecule.DropdownMenu class="mr-2" items={menu}>
                        <atom.Icon fa="fas fa-cog" />
                    </molecule.DropdownMenu>
                </div>
            </div>;
        }
    }
}