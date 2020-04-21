// Copyright (c) Jupyter Development Team.
// Distributed under the terms of the Modified BSD License.
/*-----------------------------------------------------------------------------
| Copyright (c) 2014-2017, PhosphorJS Contributors
|
| Distributed under the terms of the BSD 3-Clause License.
|
| The full license is in the file LICENSE, distributed with this software.
|----------------------------------------------------------------------------*/
import 'es6-promise/auto'; // polyfill Promise on IE

import { CommandRegistry } from '@lumino/commands';

import { Message } from '@lumino/messaging';

import {
  BoxPanel,
  DockPanel,
  Menu,
  MenuBar,
  Panel,
  SplitPanel,
  Widget,
  TabBar,
  Title,
} from '@lumino/widgets';

import {
  ElementDataset,
  ElementInlineStyle,
  VirtualDOM,
  VirtualElement,
  h,
} from '@lumino/virtualdom';

import '../style/index.css';

window.onload = () => {
  Widget.attach(new Shell(), document.body);
};

class Shell extends Panel {
  commands: CommandRegistry;
  menu: MenuBar;
  main: BoxPanel;
  split: SplitPanel;
  mainArea: DockPanel;
  leftArea: SideBar;

  constructor() {
    super();

    this.commands = new CommandRegistry();
    this.menu = this.setupMenuBar();

    this.mainArea = new DockPanel();
    this.mainArea.addWidget(new ContentWidget('Red'));

    this.leftArea = new SideBar();

    let blue = new ContentWidget('Blue', 'fa-cog');
    blue;
    let yellow = new ContentWidget('Yellow', 'fa-flask');
    this.leftArea.addWidget(blue);
    this.leftArea.addWidget(yellow);

    SplitPanel.setStretch(this.leftArea, 1);
    SplitPanel.setStretch(this.mainArea, 5);

    this.split = new SplitPanel({ spacing: 0 });
    this.split.addWidget(this.leftArea);
    this.split.addWidget(this.mainArea);
    BoxPanel.setStretch(this.split, 1);
    this.main = new BoxPanel({ spacing: 2 });
    this.main.id = 'main';
    this.main.addWidget(this.split);
    this.addWidget(this.menu);
    this.addWidget(this.split);
  }

  onAfterAttach(msg: Message) {
    window.onresize = () => {
      this.main.update();
      //todo fix on resize
      console.log('onresize');
    };
  }

  setupMenuBar(): MenuBar {
    let menu1 = this.createMenu();
    menu1.title.label = 'File';
    menu1.title.mnemonic = 0;

    let menu2 = this.createMenu();
    menu2.title.label = 'Edit';
    menu2.title.mnemonic = 0;

    let menu3 = this.createMenu();
    menu3.title.label = 'View';
    menu3.title.mnemonic = 0;

    let bar = new MenuBar();
    bar.addMenu(menu1);
    bar.addMenu(menu2);
    bar.addMenu(menu3);
    bar.id = 'menuBar';

    return bar;
  }

  createMenu(): Menu {
    let sub1 = new Menu({ commands: this.commands });
    sub1.title.label = 'More...';
    sub1.title.mnemonic = 0;
    sub1.addItem({ command: 'example:one' });
    sub1.addItem({ command: 'example:two' });
    sub1.addItem({ command: 'example:three' });
    sub1.addItem({ command: 'example:four' });

    let sub2 = new Menu({ commands: this.commands });
    sub2.title.label = 'More...';
    sub2.title.mnemonic = 0;
    sub2.addItem({ command: 'example:one' });
    sub2.addItem({ command: 'example:two' });
    sub2.addItem({ command: 'example:three' });
    sub2.addItem({ command: 'example:four' });
    sub2.addItem({ type: 'submenu', submenu: sub1 });

    let root = new Menu({ commands: this.commands });
    root.addItem({ type: 'submenu', submenu: sub2 });

    return root;
  }
}

class SideBar extends BoxPanel {
  dock: DockPanel;
  bar: TabBar<Widget>;

  constructor() {
    super({ direction: 'left-to-right' });
    this.title.closable = true;

    // dockpanel has tabs but are hidden in single-document
    this.dock = new DockPanel({ mode: 'single-document' });

    // TODO: listen for signals about dropping new tabs,
    // since they need to be added via this.addWidget to get
    // added to the bar as well

    // we setup our own vertical tab bar
    class MyTabBarRenderer extends TabBar.Renderer {
      //overwrite to hide label and close icon
      renderTab(data: TabBar.IRenderData<any>): VirtualElement {
        let title = data.title.caption;
        let key = this.createTabKey(data);
        let style = this.createTabStyle(data);
        let className = this.createTabClass(data);
        let dataset = this.createTabDataset(data);
        return h.li(
          { key, className, title, style, dataset },
          this.renderIcon(data)
        );
      }
    }

    this.bar = new TabBar<Widget>({
      orientation: 'vertical',
      tabsMovable: true,
      renderer: new MyTabBarRenderer(),
    });

    BoxPanel.setStretch(this.bar, 1);
    BoxPanel.setStretch(this.dock, 5);

    super.addWidget(this.bar);
    super.addWidget(this.dock);
  }

  addWidget(w: Widget) {
    this.dock.addWidget(w);

    this.bar.addTab(w.title);
    this.bar.tabActivateRequested.connect((tab, args) => {
      this.dock.activateWidget(args.title.owner);
    });
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

  constructor(name: string, icon?: string) {
    super({ node: ContentWidget.createNode() });
    this.setFlag(Widget.Flag.DisallowLayout);
    this.addClass('content');
    this.addClass(name.toLowerCase());
    this.title.label = name;
    this.title.closable = true;
    this.title.caption = `Long description for: ${name}`;
    if (icon) {
      this.title.iconClass = `fa ${icon} fa-2x`;
    }
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
