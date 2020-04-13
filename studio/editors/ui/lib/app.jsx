import * as form from '/views/ui/lib/form.js';
import * as atom from '/views/ui/lib/atom.js';
import * as field from '/views/ui/lib/field.js';
import * as client from '/views/ui/lib/client.js';

import * as inspector from '/object/lib/app.js';

export function App(initial) {
    return {
        view: function (vnode) {
            return <main>
                <h1>UI Demo</h1>
                <div class="flex">
                    <div class="m-2" style={{ "width": "300px" }}>
                        <inspector.ObjectHeader />
                        <field.ComponentPanel label="New Test">
                            <field.LabeledField label="Every">
                                <form.TextInput value="1" />
                                <form.SelectInput value="minutes">
                                    <option>minutes</option>
                                    <option>hours</option>
                                </form.SelectInput>
                            </field.LabeledField>
                        </field.ComponentPanel>
                        <field.ComponentPanel label="MinHour Demo">
                            <field.LabeledField label="Every">
                                <form.TextInput value="1" />
                                <form.SelectInput value="minutes">
                                    <option>minutes</option>
                                    <option>hours</option>
                                </form.SelectInput>
                            </field.LabeledField>
                            <field.Nested label="Advanced">
                                <field.Nested label="Limit Times">
                                    <field.LabeledField label="From">
                                        <form.TimeInput />
                                    </field.LabeledField>
                                    <field.LabeledField label="Until">
                                        <form.TimeInput />
                                    </field.LabeledField>
                                </field.Nested>
                                <field.Nested label="Limit Duration">
                                    <field.LabeledField label="End At">
                                        <form.DateInput />
                                    </field.LabeledField>
                                    <field.LabeledField label="End After">
                                        <form.NumberInput />
                                    </field.LabeledField>
                                    <field.LabeledField label="End on error">
                                        <form.CheckboxInput />
                                    </field.LabeledField>
                                </field.Nested>
                            </field.Nested>
                        </field.ComponentPanel>

                        <field.ComponentPanel label="Days Demo">
                            <field.LabeledField label="Period Size">
                                <form.NumberInput value="1" />
                            </field.LabeledField>
                            <field.LabeledField label="Period Unit">
                                <form.SelectInput value="days">
                                    <option>days</option>
                                </form.SelectInput>
                            </field.LabeledField>
                            <field.Collection label="Time of Day">
                                <form.TimeInput />
                            </field.Collection>
                            <field.Nested label="Advanced">
                                <field.Nested label="Limit Months of the Year">
                                    <div class="grid grid-flow-col grid-cols-2 grid-rows-6 mr-6">
                                        <form.CheckboxInput label="Jan" />
                                        <form.CheckboxInput label="Feb" />
                                        <form.CheckboxInput label="Mar" />
                                        <form.CheckboxInput label="Apr" />
                                        <form.CheckboxInput label="May" />
                                        <form.CheckboxInput label="Jun" />
                                        <form.CheckboxInput label="Jul" />
                                        <form.CheckboxInput label="Aug" />
                                        <form.CheckboxInput label="Sep" />
                                        <form.CheckboxInput label="Oct" />
                                        <form.CheckboxInput label="Nov" />
                                        <form.CheckboxInput label="Dec" />
                                    </div>
                                </field.Nested>
                                <field.Nested label="Limit Duration">
                                    <field.LabeledField label="End At">
                                        <form.DateInput />
                                    </field.LabeledField>
                                    <field.LabeledField label="End After">
                                        <form.NumberInput />
                                    </field.LabeledField>
                                    <field.LabeledField label="End on error">
                                        <form.CheckboxInput />
                                    </field.LabeledField>
                                </field.Nested>

                            </field.Nested>
                        </field.ComponentPanel>

                        <field.ComponentPanel label="Weeks Demo">
                            <field.LabeledField label="Period Size">
                                <form.NumberInput value="1" />
                            </field.LabeledField>
                            <field.LabeledField label="Period Unit">
                                <form.SelectInput value="weeks">
                                    <option>weeks</option>
                                </form.SelectInput>
                            </field.LabeledField>
                            <field.Collection label="Time of Day" expanded={true}>
                                <form.TimeInput />
                            </field.Collection>
                            <field.Nested label="Day of the Week" expanded={true}>
                                <form.CheckboxInput label="Monday" />
                                <form.CheckboxInput label="Tuesday" />
                                <form.CheckboxInput label="Wednesday" />
                                <form.CheckboxInput label="Thursday" />
                                <form.CheckboxInput label="Friday" />
                                <form.CheckboxInput label="Saturday" />
                                <form.CheckboxInput label="Sunday" />
                            </field.Nested>
                            <field.Nested label="Advanced">
                                <field.Nested label="Limit Months of the Year">
                                    <div class="grid grid-flow-col grid-cols-2 grid-rows-6 mr-6">
                                        <form.CheckboxInput label="Jan" />
                                        <form.CheckboxInput label="Feb" />
                                        <form.CheckboxInput label="Mar" />
                                        <form.CheckboxInput label="Apr" />
                                        <form.CheckboxInput label="May" />
                                        <form.CheckboxInput label="Jun" />
                                        <form.CheckboxInput label="Jul" />
                                        <form.CheckboxInput label="Aug" />
                                        <form.CheckboxInput label="Sep" />
                                        <form.CheckboxInput label="Oct" />
                                        <form.CheckboxInput label="Nov" />
                                        <form.CheckboxInput label="Dec" />
                                    </div>
                                </field.Nested>
                                <field.Nested label="Limit Duration">
                                    <field.LabeledField label="End At">
                                        <form.DateInput />
                                    </field.LabeledField>
                                    <field.LabeledField label="End After">
                                        <form.NumberInput />
                                    </field.LabeledField>
                                    <field.LabeledField label="End on error">
                                        <form.CheckboxInput />
                                    </field.LabeledField>
                                </field.Nested>

                            </field.Nested>
                        </field.ComponentPanel>

                        <field.ComponentPanel label="Month Demo">
                            <field.LabeledField label="Period Size">
                                <form.NumberInput value="1" />
                            </field.LabeledField>
                            <field.LabeledField label="Period Unit">
                                <form.SelectInput value="months">
                                    <option>months</option>
                                </form.SelectInput>
                            </field.LabeledField>
                            <field.LabeledField label="Day Selection">
                                <form.SelectInput>
                                    <option>day of month</option>
                                    <option>day of week</option>
                                </form.SelectInput>
                            </field.LabeledField>
                            <field.Collection label="Time of Day" expanded={true}>
                                <form.TimeInput />
                            </field.Collection>
                            <field.Collection label="Day of the Month" expanded={true}>
                                <form.NumberInput value="7" />
                            </field.Collection>
                            <field.Collection label="Week of the Month" expanded={true}>
                                <form.NumberInput value="1" />
                            </field.Collection>
                            <field.Nested label="Day of the Week" expanded={true}>
                                <form.CheckboxInput label="Monday" />
                                <form.CheckboxInput label="Tuesday" />
                                <form.CheckboxInput label="Wednesday" />
                                <form.CheckboxInput label="Thursday" />
                                <form.CheckboxInput label="Friday" />
                                <form.CheckboxInput label="Saturday" />
                                <form.CheckboxInput label="Sunday" />
                            </field.Nested>
                            <field.Nested label="Advanced">
                                <field.Nested label="Limit Months of the Year">
                                    <div class="grid grid-flow-col grid-cols-2 grid-rows-6 mr-6">
                                        <form.CheckboxInput label="Jan" />
                                        <form.CheckboxInput label="Feb" />
                                        <form.CheckboxInput label="Mar" />
                                        <form.CheckboxInput label="Apr" />
                                        <form.CheckboxInput label="May" />
                                        <form.CheckboxInput label="Jun" />
                                        <form.CheckboxInput label="Jul" />
                                        <form.CheckboxInput label="Aug" />
                                        <form.CheckboxInput label="Sep" />
                                        <form.CheckboxInput label="Oct" />
                                        <form.CheckboxInput label="Nov" />
                                        <form.CheckboxInput label="Dec" />
                                    </div>
                                </field.Nested>
                                <field.Nested label="Limit Duration">
                                    <field.LabeledField label="End At">
                                        <form.DateInput />
                                    </field.LabeledField>
                                    <field.LabeledField label="End After">
                                        <form.NumberInput />
                                    </field.LabeledField>
                                    <field.LabeledField label="End on error">
                                        <form.CheckboxInput />
                                    </field.LabeledField>
                                </field.Nested>

                            </field.Nested>
                        </field.ComponentPanel>

                    </div>
                    <ul class="m-2 w-64">
                        <li class="m-2 inline-flex w-full">
                            <div><form.CheckboxInput /></div>
                            <div class="flex-grow"><form.SliderInput value="50" step="5" /></div>
                            <div class="mr-4"><form.KnobInput /></div>
                        </li>
                        <li class="m-2 flex">
                            <div class="flex-grow align-center"><atom.Indicator active={false} /></div>
                            <div class="flex-grow"><atom.Indicator active color="#00ff00" /></div>
                            <div class="flex-grow"><atom.Indicator active color="yellow" /></div>
                            <div><atom.Indicator active color="red" /></div>
                        </li>
                        <li class="m-2"><atom.Button /></li>
                        <li class="m-2"><form.TextInput value="Hello" /></li>
                        <li class="m-2"><form.PasswordInput value="Hello" /></li>
                        <li class="m-2"><form.NumberInput value="100" /></li>
                        <li class="m-2">
                            <form.SelectInput value="Bar">
                                <option>Foo</option>
                                <option>Baz</option>
                                <option>Bar</option>
                            </form.SelectInput>
                        </li>
                        <li class="m-2"><form.ReferenceInput placeholder="os.FileInfo" /></li>
                        <li class="m-2"><form.DateInput /></li>
                        <li class="m-2"><form.TimeInput /></li>
                        <li class="m-2"><form.ColorInput /></li>
                        <li class="m-2">
                            <field.ComponentPanel label="http.Server">
                                <form.DateInput />
                                <form.DateInput />
                                <form.DateInput />
                                <field.Collection label="Stuff">
                                    <form.ColorInput />
                                    <form.ColorInput />
                                </field.Collection>
                                <field.Collection label="Stuff2">
                                    <form.ColorInput />
                                    <form.ColorInput />
                                </field.Collection>
                            </field.ComponentPanel>
                        </li>
                        <li class="m-2">
                            <field.Collection label="Container">
                                <form.ColorInput />
                                <form.DateInput />
                                <form.NumberInput />
                                <form.TimeInput />
                                <form.SelectInput />
                                <form.ReferenceInput />
                                <form.TextInput />
                            </field.Collection>
                        </li>
                    </ul>


                </div>
            </main>;
        }
    }
}

// not canonical, pls delete
export function Inspector(initial) {
    let remote = { components: [] };
    var node;
    let path = "/Users/progrium/Projects/tractor/local/workspace";
    let session = new client.Session(path, (client) => {
        client.call("subscribe");
    });
    session.api.handle("shutdown", {
        "serveRPC": async (r, c) => {
            console.log("DEBUG: reload/shutdown received...");
            session.reconnect();
            r.return();
        }
    });
    session.api.handle("state", {
        "serveRPC": async (r, c) => {
            remote = await c.decode();
            node = remote.nodes[remote.selectedNode];
            // if (data.selectedNode && !this.state.lastSelected) {
            //     this.setState({"lastSelected": data.selectedNode});
            // }
            m.redraw();
            r.return();
        }
    });

    return {
        view: function (vnode) {
            if (node) {
                return <section>
                    <field.ObjectHeader />
                    {node.components.map((el) =>
                        <field.ComponentPanel label={el.name}>
                            {(el.fields || []).map((f, idx) => // this is where you can customize ui
                                <field.ComponentField key={idx} field={f} />
                            )}
                        </field.ComponentPanel>
                    )}
                </section>;
            }
            return <section>no node selected</section>;
        }
    }
}
