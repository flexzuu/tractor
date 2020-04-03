import { ContainerModule } from 'inversify';
import { EditorViewWidget, EditorViewOptions } from './editorview-widget';
import { EditorViewContribution } from './editorview-contribution';
import { bindViewContribution, FrontendApplicationContribution, WidgetFactory } from '@theia/core/lib/browser';

import '../../src/browser/style/index.css';

export default new ContainerModule(bind => {
    bindViewContribution(bind, EditorViewContribution);
    bind(FrontendApplicationContribution).toService(EditorViewContribution);
    bind(EditorViewWidget).toSelf();
    bind(WidgetFactory).toDynamicValue(ctx => ({
        id: EditorViewWidget.FACTORY_ID,
        createWidget: async (options: EditorViewOptions) => {
            const child = ctx.container.createChild();
            child.bind(EditorViewOptions).toConstantValue(options);
            return child.get(EditorViewWidget);
        }
        //createWidget: () => ctx.container.get<EditorViewWidget>(EditorViewWidget)
    })).inSingletonScope();
});
