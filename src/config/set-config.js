import {SetService} from "../set/set-service.js";
import {SetDao} from "../set/data/set-dao.js";
import {SetViewResource} from "../set/set-view-resource.js";

export function configureSetObjets(dbClientFactory) {
    const setDao = new SetDao();
    const setService = new SetService(dbClientFactory, setDao);
    const setViewResource = new SetViewResource(setService);

    return {
        setDao,
        setService,
        setViewResource
    };
}
