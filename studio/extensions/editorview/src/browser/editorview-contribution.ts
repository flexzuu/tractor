import { injectable, inject } from 'inversify';
import { MenuModelRegistry, ILogger } from '@theia/core';
import { EditorViewWidget } from './editorview-widget';
import { EditorTypes } from './editorviews';
import { CommonMenus } from '@theia/core/lib/browser/common-frontend-contribution';
import { AbstractViewContribution } from '@theia/core/lib/browser';
import { CommandRegistry } from '@theia/core/lib/common/command';
import { WidgetManager } from '@theia/core/lib/browser/widget-manager';
import { ApplicationShell } from '@theia/core/lib/browser/shell/application-shell';
import { OpenerService } from '@theia/core/lib/browser/opener-service';
import URI from '@theia/core/lib/common/uri';

@injectable()
export class EditorViewContribution extends AbstractViewContribution<EditorViewWidget> {

    @inject(WidgetManager)
    protected readonly widgets: WidgetManager;

    @inject(ApplicationShell)
    protected readonly shell: ApplicationShell;

    @inject(ILogger)
    protected readonly logger: ILogger;

    @inject(OpenerService)
    protected readonly openerService: OpenerService;

    /**
     * `AbstractViewContribution` handles the creation and registering
     *  of the widget including commands, menus, and keybindings.
     * 
     * We can pass `defaultWidgetOptions` which define widget properties such as 
     * its location `area` (`main`, `left`, `right`, `bottom`), `mode`, and `ref`.
     * 
     */
    constructor() {
        super({
            widgetId: EditorViewWidget.FACTORY_ID,
            widgetName: "EditorView",
            defaultWidgetOptions: { area: 'main' },
        });
    }


    onStart(app: any): void {
        window.addEventListener("message", async (e) => {
            let message = e.data;
            if (message.event) {
                switch (message.event) {
                    case 'edit':
                        if (message.path !== undefined) {
                            const uriArg = new URI(`file://${message.path}`);
                            const opener = await this.openerService.getOpener(uriArg, {});
                            await opener.open(uriArg, {});
                            return;
                        }
                        return;

                }

            }
        });
    }

    /**
     * Example command registration to open the widget from the menu, and quick-open.
     * For a simpler use case, it is possible to simply call:
     ```ts
        super.registerCommands(commands)
     ```
     *
     * For more flexibility, we can pass `OpenViewArguments` which define 
     * options on how to handle opening the widget:
     * 
     ```ts
        toggle?: boolean
        activate?: boolean;
        reveal?: boolean;
     ```
     *
     * @param commands
     */
    registerCommands(commands: CommandRegistry): void {
        EditorTypes.forEach((editor) => {
            commands.registerCommand({ id: 'editor:' + editor.name, label: "Show: " + editor.label }, {
                execute: async () => {
                    const widget = await this.widgets.getOrCreateWidget<EditorViewWidget>("editorview", editor);
                    this.shell.addWidget(widget, <ApplicationShell.WidgetOptions>{ 'area': editor.area });
                    this.shell.revealWidget(widget.id);
                }
            });
        })
    }

    /**
     * Example menu registration to contribute a menu item used to open the widget.
     * Default location when extending the `AbstractViewContribution` is the `View` main-menu item.
     * 
     * We can however define new menu path locations in the following way:
     ```ts
        menus.registerMenuAction(CommonMenus.HELP, {
            commandId: 'id',
            label: 'label'
        });
     ```
     * 
     * @param menus
     */
    registerMenus(menus: MenuModelRegistry): void {
        EditorTypes.forEach((editor) => {
            menus.registerMenuAction(CommonMenus.VIEW_VIEWS, {
                commandId: 'editor:' + editor.name,
                label: editor.label
            });
        });

    }
}
