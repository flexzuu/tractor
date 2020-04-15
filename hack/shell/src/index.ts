// Copyright (c) Jupyter Development Team.
// Distributed under the terms of the Modified BSD License.
/*-----------------------------------------------------------------------------
| Copyright (c) 2014-2017, PhosphorJS Contributors
|
| Distributed under the terms of the BSD 3-Clause License.
|
| The full license is in the file LICENSE, distributed with this software.
|----------------------------------------------------------------------------*/
import 'es6-promise/auto';  // polyfill Promise on IE

import {
    CommandRegistry
} from '@lumino/commands';

import {
    Message
} from '@lumino/messaging';

import {
    BoxPanel, SplitPanel, BoxLayout, DockPanel, Menu, MenuBar, Widget, TabBar, StackedPanel, TabPanel
} from '@lumino/widgets';

import '../style/index.css';


const commands = new CommandRegistry();



class SideBar extends BoxPanel {

    dock: DockPanel;
    bar: BoxPanel;

    constructor() {
        super({ "direction": "left-to-right" });
        this.title.closable = true;

        // dockpanel has tabs but are hidden in single-document
        this.dock = new DockPanel({ 'mode': 'single-document' });

        // TODO: listen for signals about dropping new tabs, 
        // since they need to be added via this.addWidget to get
        // added to the bar as well

        // bar visualizes and drives the tabs for this.dock
        this.bar = new BoxPanel();


        BoxPanel.setStretch(this.bar, 1);
        BoxPanel.setStretch(this.dock, 5);

        super.addWidget(this.bar);
        super.addWidget(this.dock);
    }

    addWidget(w: Widget) {
        this.dock.addWidget(w);
        let icon = new IconWidget("fa fa-cog fa-2x");
        this.bar.addWidget(icon);
    }

}


class IconWidget extends Widget {

    static createNode(name: string): HTMLElement {
        let node = document.createElement('div');
        node.className = name;
        return node;
    }

    constructor(name: string) {
        super({ node: IconWidget.createNode(name) });
        this.title.label = name;
    }

}


class ContentWidget extends Widget {

    static createNode(): HTMLElement {
        let node = document.createElement('div');
        let content = document.createElement('div');
        let input = document.createElement('input');
        input.placeholder = 'Placeholder...';
        content.appendChild(input);
        node.appendChild(content);
        return node;
    }

    constructor(name: string) {
        super({ node: ContentWidget.createNode() });
        this.setFlag(Widget.Flag.DisallowLayout);
        this.addClass('content');
        this.addClass(name.toLowerCase());
        this.title.label = name;
        this.title.closable = true;
        this.title.caption = `Long description for: ${name}`;
    }

    get inputNode(): HTMLInputElement {
        return this.node.getElementsByTagName('input')[0] as HTMLInputElement;
    }

    protected onActivateRequest(msg: Message): void {
        if (this.isAttached) {
            this.inputNode.focus();
        }
    }
}


function main(): void {

    let bar = setupMenuBar();

    let r1 = new ContentWidget('Red');
    let b1 = new ContentWidget('Blue');
    let g1 = new ContentWidget('Green');
    let y1 = new ContentWidget('Yellow');


    let mainArea = new DockPanel();
    mainArea.addWidget(r1);

    let leftArea = new SideBar();
    leftArea.addWidget(b1);
    leftArea.addWidget(g1);

    SplitPanel.setStretch(leftArea, 1);
    SplitPanel.setStretch(mainArea, 5);
    // SplitPanel.setStretch(rightArea, 1);

    let main = new SplitPanel({ spacing: 0 });
    main.id = 'main';
    main.addWidget(leftArea);
    main.addWidget(mainArea);
    // main.addWidget(rightArea);

    window.onresize = () => { main.update(); };

    Widget.attach(bar, document.body);
    Widget.attach(main, document.body);
}


window.onload = main;


function createMenu(): Menu {
    let sub1 = new Menu({ commands });
    sub1.title.label = 'More...';
    sub1.title.mnemonic = 0;
    sub1.addItem({ command: 'example:one' });
    sub1.addItem({ command: 'example:two' });
    sub1.addItem({ command: 'example:three' });
    sub1.addItem({ command: 'example:four' });

    let sub2 = new Menu({ commands });
    sub2.title.label = 'More...';
    sub2.title.mnemonic = 0;
    sub2.addItem({ command: 'example:one' });
    sub2.addItem({ command: 'example:two' });
    sub2.addItem({ command: 'example:three' });
    sub2.addItem({ command: 'example:four' });
    sub2.addItem({ type: 'submenu', submenu: sub1 });

    let root = new Menu({ commands });
    root.addItem({ type: 'submenu', submenu: sub2 });

    return root;
}

function setupMenuBar() {

    let menu1 = createMenu();
    menu1.title.label = 'File';
    menu1.title.mnemonic = 0;

    let menu2 = createMenu();
    menu2.title.label = 'Edit';
    menu2.title.mnemonic = 0;

    let menu3 = createMenu();
    menu3.title.label = 'View';
    menu3.title.mnemonic = 0;

    let bar = new MenuBar();
    bar.addMenu(menu1);
    bar.addMenu(menu2);
    bar.addMenu(menu3);
    bar.id = 'menuBar';

    return bar;
}
