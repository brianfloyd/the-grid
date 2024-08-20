import {SetService} from "../set/set-service.js";
import {SetDao} from "../set/data/set-dao.js";

export function configureSetObjets(dbClientFactory) {
    const setDao = new SetDao();
    const setService = new SetService(dbClientFactory, setDao);

    return {
        setDao,
        setService
    };
}
