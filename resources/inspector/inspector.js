const { Column, Checkbox, Input, Message, Delete, Box, Heading, Content, List, Dropdown, Breadcrumb, Button, Icon, Level } = rbx;

function remoteAction(action, args) {
    switch (action) {
        case "setValue":
        case "setExpression":
        case "callMethod":
        case "removeComponent":
        case "appendComponent":
        case "addDelegate":
            console.log(action);
            window.rpc.call(action, args)
        default:
            throw "unknown action";
    }
    
}

function FieldControl(props) {
    const [exprMode, setExprMode] = React.useState(false);
    let onChange = (event) => remoteAction("setValue", { "Path": props.path, "Value": event.target.value });
    let readOnly = (props.expression || "").length > 0;
    function typedControl() {
        if (exprMode) {
            onChange = (event) => remoteAction("setExpression", { "Path": props.path, "Value": event.target.value });
            return <Input type="text" size="small" style={{ height: "22px", color: "white", backgroundColor: "#555", fontFamily: "monospace" }} onChange={onChange} value={props.expression} />
        }
        switch (props.type) {
            case "string":
                return <Input type="text" readOnly={readOnly} size="small" onChange={onChange} value={props.value} />
            case "boolean":
                onChange = (event) => remoteAction("setValue", { "Path": this.props.path, "Value": event.target.checked });
                return <Checkbox onChange={onChange} readOnly={readOnly} checked={props.value} />
            
            case "number":
                onChange = (val) => remoteAction("setValue", { "Path": this.props.path, "IntValue": val });
                return <Input type="number" readOnly={readOnly} style={{ width: "100%" }} size="small" onChange={onChange} value={props.value} />
            default:
                if (props.type.startsWith("reference:")) {
                    var refType = props.type.split(":")[1];
                    onChange = (val, label, extra) => window.rpc.call("setValue", { "Path": this.props.path, "RefValue": val + "/" + refType });
                    const data = ["one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten"];
                    // TODO: finish me
                    return <Autocomplete data={data} value={props.value} />;
                } else {
                    return "???";
                }
        }
    }
    return (
        <div className="field is-grouped">
            <div className="control is-expanded" style={{marginRight: "0px", height: "22px", }}>
                {typedControl()}
            </div>
            <div className="control">
                <div className="image is-16x16" style={{margin: "3px"}}>
                    <img src={window.baseUri+"/inspector/function-icon.png"} onClick={() => setExprMode(!exprMode)} />
                </div>
            </div>
        </div>
    );
}

function FieldRow(props) {
    const children = props.children.slice(0);
    const label = children.shift();
    return (
        <Column.Group style={{display: "flex", marginBottom: "0"}}>
            <Column style={{ fontSize: "smaller"}}>
                {label}
            </Column>
            <Column style={{flexGrow: "2"}}>
                {children}
            </Column>
        </Column.Group>
    );
}

function LabeledField(props) {
    return (
        <FieldRow>
            <span>{props.label}</span>
            {props.children}
        </FieldRow>
    );
}

function KeyedField(props) {
    return (
        <FieldRow>
            <Input type="text" value={props.name} />
            {props.children}
        </FieldRow>
    );
  }

function EmbeddedFields(props) {
    const fields = props.fields || [];
    const [open, setOpen] = React.useState(false);
    return (
        <div>
            <header onClick={() => setOpen(!open)}>
                <Arrow opened={open} />
                <span style={{fontSize: "smaller"}}>{props.name}</span>
            </header>
            {open &&
                <Content style={{marginLeft: "20px", marginTop: "8px"}}>
                    {props.children}
                </Content>
            }
        </div>
    );
}


function ComponentField(props) {
    switch (props.type) {
        case "boolean":
        case "string":
        case "number":
        case "reference":
            return <LabeledField key={props.eventKey} label={props.name}><FieldControl {...props} /></LabeledField>
        case "struct":
        case "map":
        case "array":
            let fields = props.fields || [];
            if (props.type == "array") {
                fields = fields.map((obj, idx) => { obj.name = "Element "+idx; return obj; })
            }
            return (
                <EmbeddedFields {...props}>
                    {(props.type == "array") &&
                        <LabeledField key="count" label="Count">
                            <Input type="number" style={{ width: "100%" }} size="small" value={fields.length} />
                        </LabeledField>
                    }
                    {fields.map((field) => {
                        if (props.type == "map") {
                            return <KeyedField key={field.name} name={field.name}><FieldControl {...field} /></KeyedField>;
                        } else {
                            return <ComponentField key={field.name} {...field} />;
                        }
                    })}
                    {props.children}
                </EmbeddedFields>
            );
        default:
            return <LabeledField key={props.eventKey} label={props.name}><FieldControl {...props} /></LabeledField>
    }
}

function ComponentFields(props) {
    const addKey = (obj, key) => { obj.key = key; return obj};
    return (props.fields || []).map((el, idx) => 
        <ComponentField {...addKey(el, idx)} />
    );
}

function ComponentManageMenu(props) {
    return (
        <Dropdown align="right" {...props}>
            <Dropdown.Trigger>
                <Icon size="small"><i className="fas fa-cog"></i></Icon>
            </Dropdown.Trigger>
            <Dropdown.Menu>
                <Dropdown.Content>
                    <Dropdown.Item onClick={() => { /* TODO */}}>Edit</Dropdown.Item>
                    <Dropdown.Item onClick={() => remoteAction("removeComponent", props.component.path)}>Delete</Dropdown.Item>
                </Dropdown.Content>
            </Dropdown.Menu>
        </Dropdown>
    );
}


function ComponentInspector(props) {
    const [open, setOpen] = React.useState(false);
    const headingStyle = {
        marginBottom: "0px", 
        overflow: "auto", 
        paddingBottom: open ? "15px" : "0"
    };
    return (
        <List.Item as="div">
            <ComponentManageMenu component={props.component} style={{float: "right"}} />
            <Heading as="div" onClick={() => setOpen(!open)} style={headingStyle}>
                <Arrow opened={open} />
                <span>{props.component.name}</span>
            </Heading>
            {open &&
                <Content>
                    <ComponentFields fields={props.component.fields} />
                    <ComponentButtons buttons={props.component.buttons} />
                </Content>
            }
        </List.Item>
    );
}

function ComponentButtons(props) {
    const buttonStyle = {marginTop: "10px", width: "100%"};
    function onClicker(button) {
        if (button.onclick !== "") {
            return (event) => eval(button.onclick); 
        } else {
            return (event) => remoteAction("callMethod", event.target.value);
        }
    }
    return (props.buttons||[]).map((button, idx) =>
        <Button size="small" onClick={onClicker(button)} key={idx} value={button.path} style={buttonStyle}>
            {button.name}
        </Button>
    );
}

function AddComponentButton(props) {
    const [open, setOpen] = React.useState(false);
    const [filterFocus, setFilterFocus] = React.useState(false);
    
    const fakeData = ["one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten"];
    const dropdownRef = React.useRef(null);
    const filterRef = React.useRef(null);
    const onBlur = () => {
        setTimeout(() => {
            if (dropdownRef.current.contains(document.activeElement)) {
                return;
            }
            setOpen(false);
        }, 200);
    }
    React.useEffect(() => filterRef.current.focus());
    function onClicker(component) {
        return () => remoteAction("appendComponent", component);
    }
    return (
        <section style={{textAlign: "center", marginBottom: "20px"}}>
            <Dropdown up managed active={open} onBlur={onBlur} ref={dropdownRef}>
                <Dropdown.Trigger>
                    <Button onClick={() => setOpen(!open)}>
                        <span>Add Component</span>
                        <Icon size="small"><i className="fas fa-angle-up"></i></Icon>
                    </Button>
                </Dropdown.Trigger>
                <Dropdown.Menu>
                    <Dropdown.Content>
                        <Dropdown.Item as="div">
                            <Input type="text" size="small" ref={filterRef} />
                        </Dropdown.Item>
                        <Dropdown.Item as="div" style={{maxHeight: "100px", overflowY: "scroll", textAlign: "left"}}>
                            {fakeData.map((item, idx) =>
                                <div key={idx} onClick={onClicker(item)}>{item}</div>
                            )}
                        </Dropdown.Item>
                    </Dropdown.Content>
                </Dropdown.Menu>
            </Dropdown>
        </section>
    );
}

function Inspector(props) {
    const node = props.node || {name: "no node", components: []};
    return (
        <section>
            <Breadcrumb style={{marginBottom: "0"}}>
                <Breadcrumb.Item as="div" style={{color: "white", margin: "5px"}}>Foobar</Breadcrumb.Item>
                <Breadcrumb.Item as="div" style={{color: "white", margin: "5px"}} active>{node.name}</Breadcrumb.Item>
            </Breadcrumb>
            <List>
                {node.components.map((com, idx) => 
                    <ComponentInspector component={com} key={"com-"+idx} />
                )}
            </List>
            <AddComponentButton />
        </section>
    );
}

function Arrow(props) {
    return (
        <Icon size="small"><i className={props.opened ? "fas fa-angle-down" : "fas fa-angle-right"}></i></Icon>
    );
}

function Autocomplete(props) {
    const [selected, setSelected] = React.useState(props.value);
    const [active, setActive] = React.useState(false);
    let filteredData = props.data || [];
    if (selected !== undefined && selected.length > 0) {
        filteredData = filteredData.filter((el) => el.startsWith(selected));
    }
    const items = filteredData.map((el) => 
        <Dropdown.Item key={el} onClick={() => {setSelected(el); setActive(false)}}>{el}</Dropdown.Item>
    )
    const onChange = (e) => {
        setSelected(e.target.value);
        setActive(true);
    };
    const onBlur = () => {
        setTimeout(() => setActive(false), 100);
    }
    return (
        <div>
            <Input type="text" value={selected} size="small" onChange={onChange} onBlur={onBlur} />
            <Dropdown.Content style={{display: active ? "block" : "none", "position": "absolute", "width": "90%"}}>
                {items}
            </Dropdown.Content>
        </div>
    );
}