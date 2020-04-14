import * as React from 'react';
import { injectable, postConstruct, inject } from 'inversify';
import { ReactWidget } from '@theia/core/lib/browser/widgets/react-widget';
import { TractorService } from 'tractor/lib/browser/tractor-service';
import { WorkspaceService } from '@theia/workspace/lib/browser';
import URI from '@theia/core/lib/common/uri';

let counter = 0;

@injectable()
export class EditorViewOptions {
    name: string;
    label: string;
    iconClass: string;
    area: string;
}

@injectable()
export class EditorViewWidget extends ReactWidget {

    static readonly FACTORY_ID = 'editorview';

    @inject(WorkspaceService)
    protected readonly workspace: WorkspaceService;

    @inject(EditorViewOptions)
    readonly options: EditorViewOptions;

    @inject(TractorService)
    readonly tractor: TractorService;

    @postConstruct()
    protected async init(): Promise<void> {
        counter++;
        this.id = EditorViewWidget.FACTORY_ID + ":" + counter; // TODO: something smarter?
        this.title.label = this.options.label;
        this.title.caption = this.options.label;
        this.title.iconClass = this.options.iconClass;
        this.title.closable = true;
        this.update();
    }

    protected render(): React.ReactNode {
        let workspacePath = new URI(this.workspace.workspace.uri).path.toString();
        let buf = Buffer.from(JSON.stringify({ "workspace": workspacePath }));
        return <iframe style={{ border: "0px", width: "100%", height: "100%" }}
            src={`http://${window.location.host}/views/${this.options.name}#${buf.toString('base64')}`}>
        </iframe>;
    }


}
