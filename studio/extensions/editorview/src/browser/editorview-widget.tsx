import * as React from 'react';
import { injectable, postConstruct, inject } from 'inversify';
import { ReactWidget } from '@theia/core/lib/browser/widgets/react-widget';
import { TractorService } from 'tractor/lib/browser/tractor-service';

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

    @inject(EditorViewOptions)
    readonly options: EditorViewOptions;

    @inject(TractorService)
    readonly tractor: TractorService;

    @postConstruct()
    protected async init(): Promise < void> {
        counter++; 
        this.id = EditorViewWidget.FACTORY_ID+":"+counter; // TODO: something smarter?
        this.title.label = this.options.label;
        this.title.caption = this.options.label;
        this.title.iconClass = this.options.iconClass;
        this.title.closable = true;
        this.update();
    }

    protected render(): React.ReactNode {
        return <iframe style={{border: "0px", width: "100%", height: "100%"}}
            src={`http://${this.tractor.editorsEndpoint}/${this.options.name}`}>
        </iframe>;
    }


}
