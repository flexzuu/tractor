/********************************************************************************
 * Copyright (C) 2017 TypeFox and others.
 *
 * This program and the accompanying materials are made available under the
 * terms of the Eclipse Public License v. 2.0 which is available at
 * http://www.eclipse.org/legal/epl-2.0.
 *
 * This Source Code may also be made available under the following Secondary
 * Licenses when the conditions for such availability set forth in the Eclipse
 * Public License v. 2.0 are satisfied: GNU General Public License, version 2
 * with the GNU Classpath Exception which is available at
 * https://www.gnu.org/software/classpath/license.html.
 *
 * SPDX-License-Identifier: EPL-2.0 OR GPL-2.0 WITH Classpath-exception-2.0
 ********************************************************************************/

import { injectable, inject } from 'inversify';
import URI from '@theia/core/lib/common/uri';
import { Event, Emitter, DisposableCollection, MenuModelRegistry } from '@theia/core';
import { WidgetFactory, TreeSelectionService } from '@theia/core/lib/browser';
import { Command, CommandRegistry } from '@theia/core/lib/common/command';
import { Widget } from '@phosphor/widgets';
import { WorkspaceService } from '@theia/workspace/lib/browser';
import { ILogger, MessageService } from '@theia/core';

import { TractorTreeWidget, ObjectNode, TractorTreeWidgetFactory } from './tractor-tree-widget';
import { TractorContextMenu, TRACTOR_CONTEXT_MENU } from './tractor-contribution';

import * as qmux from 'qmux/dist/browser/qmux.min.js';
import * as qrpc from 'qrpc';

const RetryInterval = 500;

function scheduleRetry(fn: any) {
	setTimeout(fn, RetryInterval);
}


@injectable()
export class TractorService implements WidgetFactory {

    id = TractorTreeWidget.ID;

    @inject(WorkspaceService)
    protected readonly workspace: WorkspaceService;

    @inject(MessageService)
    protected readonly messages: MessageService;

    @inject(CommandRegistry)
    protected readonly commands: CommandRegistry;

    @inject(MenuModelRegistry)
    protected readonly menus: MenuModelRegistry;

    // @inject(TreeSelectionService) 
    // protected readonly selection: TreeSelectionService;

    @inject(ILogger)
    protected readonly logger: ILogger;

    protected client: qrpc.Client;
    protected api: qrpc.API;

    public components: any[];
    public prefabs: any[];

    protected widget?: TractorTreeWidget;
    protected readonly onDidChangeEmitter = new Emitter<ObjectNode[]>();
    protected readonly onDidChangeOpenStateEmitter = new Emitter<boolean>();
    protected readonly onDidSelectEmitter = new Emitter<ObjectNode>();
    protected readonly onDidOpenEmitter = new Emitter<ObjectNode>();

    constructor(@inject(TractorTreeWidgetFactory) protected factory: TractorTreeWidgetFactory) { }

    get onDidSelect(): Event<ObjectNode> {
        return this.onDidSelectEmitter.event;
    }

    get onDidOpen(): Event<ObjectNode> {
        return this.onDidOpenEmitter.event;
    }

    get onDidChange(): Event<ObjectNode[]> {
        return this.onDidChangeEmitter.event;
    }

    get onDidChangeOpenState(): Event<boolean> {
        return this.onDidChangeOpenStateEmitter.event;
    }

    get open(): boolean {
        return this.widget !== undefined && this.widget.isVisible;
    }

    async connectAgent() {
        try {
			var conn = await qmux.DialWebsocket("ws://localhost:3001/");
		} catch (e) {
            this.logger.warn(e);
			scheduleRetry(() => this.connectAgent());
			return;
		}
        var session = new qmux.Session(conn);
        var client = new qrpc.Client(session);
        var path = new URI(this.workspace.workspace.uri).path.toString()
        var resp = await client.call("connect", path);
        this.connectWorkspace(resp.reply);
    }

    async connectWorkspace(socketPath: string) {
        this.logger.warn("attempting connect");
		try {
			var conn = await qmux.DialWebsocket("ws://localhost:3001"+socketPath);
		} catch (e) {
            this.logger.warn(e);
			scheduleRetry(() => this.connectWorkspace(socketPath));
			return;
		}
        var session = new qmux.Session(conn);
        this.api = new qrpc.API();
		this.client = new qrpc.Client(session, this.api);
		this.api.handle("shutdown", {
			"serveRPC": async (r, c) => {
                scheduleRetry(() => this.connectWorkspace(socketPath));
                r.return();
			}
        });
        this.api.handle("state", {
			"serveRPC": async (r, c) => {
                var data = await c.decode();
                //this.logger.warn(data);
                this.components = data.components;
                this.prefabs = data.prefabs;
                this.refreshRegistries();
                if (this.widget) {
                    this.widget.setData(data);
                    this.onDidChangeEmitter.fire(this.widget.rootObjects());
                }
                r.return();
			}
        });
        this.client.serveAPI();
        if (this.widget) {
            this.widget.model.onSelectionChanged(event => {
                const node = this.widget.model.selectedNodes[0];
                this.buildContextMenus(node as ObjectNode);
                this.client.call("selectNode", node.id);
            });
        }
		await this.client.call("subscribe");
    }

    buildContextMenus(node: ObjectNode) {
        const index = TractorContextMenu.COMPONENTS.length - 1;
        const menuId = TractorContextMenu.COMPONENTS[index];
        this.components.forEach((c) => {
            let id = `tractor:component-add:${c.name}`
            this.menus.unregisterMenuAction(id, TractorContextMenu.COMPONENTS);
        });
        this.menus.unregisterMenuAction('related_com', TractorContextMenu.WORKSPACE);
        this.menus.registerSubmenu(TractorContextMenu.COMPONENTS, 'Add Related');

        if (node.relatedComponents) {
            node.relatedComponents.forEach((name) => {
                let cmdId = `tractor:component-add:${name}`
                this.menus.registerMenuAction(TractorContextMenu.COMPONENTS, {
                    commandId: cmdId,
                    label: name
                });
            });
        }
    }
    
    refreshRegistries() {
        this.components.forEach((c) => {
            let id = `tractor:component-add:${c.name}`
            let label = `Add Component: ${c.name}`

            this.commands.unregisterCommand(id);
            this.commands.registerCommand({id:id, label:label}, {
                execute: () =>  {
                    let node = this.widget.model.selectedNodes[0];
                    if (node) {
                        this.addComponent(c.name, node.id);
                    }
                }
            });
            
        });
        this.prefabs.forEach((p) => {
            let id = `tractor:prefab-add:${p.id}`
            let label = `Load Prefab: ${p.name}`

            this.commands.unregisterCommand(id);
            this.commands.registerCommand({id:id, label:label}, {
                execute: () =>  {
                    let node = this.widget.model.selectedNodes[0];
                    if (node) {
                        this.loadPrefab(p.id, node.id);
                    }
                }
            });
        });
    }

    renameNode(id: string, name: string) {
		this.client.call("updateNode", {
			"ID": id,
			"Name": name
		});
	}

    addNode(name: string, parentId?: string) {
		this.client.call("appendNode", {"ID": parentId||"", "Name": name});
	}

	deleteNode(id: string) {
		this.client.call("deleteNode", id);
    }

    addComponent(component: string, nodeId: string) {
        this.client.call("appendComponent", {ID: nodeId, Name: component});
    }

    loadPrefab(id: string, nodeId: string) {
        this.client.call("loadPrefab", {ID: nodeId, Name: id});
    }
    

    createWidget(): Promise<Widget> {
        this.widget = this.factory();
        const disposables = new DisposableCollection();
        disposables.push(this.widget.onDidChangeOpenStateEmitter.event(open => this.onDidChangeOpenStateEmitter.fire(open)));
        disposables.push(this.widget.model.onOpenNode(node => this.onDidOpenEmitter.fire(node as ObjectNode)));
        disposables.push(this.widget.model.onSelectionChanged(selection => this.onDidSelectEmitter.fire(selection[0] as ObjectNode)));
        this.widget.disposed.connect(() => {
            this.widget = undefined;
            disposables.dispose();
        });
        return Promise.resolve(this.widget);
    }
}