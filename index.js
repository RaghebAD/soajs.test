'use strict';
const soajs = require('soajs');
const config = require('./soa.json');
config.packagejson = require("./package.json");
const service = new soajs.server.service(config);
const coreModules = require("soajs");
const provision = coreModules.provision;
const express = require('express');

const sApp = express();
const mApp = express();

function startServer(serverConfig, callback) {
	let sReply = {
		'result': true,
		'data': {
			'service': serverConfig.name,
			'type': 'rest',
			'route': "/"
		}
	};
    let mReply = {
        'result': true,
        'ts': Date.now(),
        'service': {
            'service': serverConfig.name,
            'type': 'rest',
            'route': "/heartbeat"
        }
    };

    sApp.get('/', (req, res) => res.json(sReply));
    mApp.get('/heartbeat', (req, res) => res.json(mReply));
    mApp.get('/maintenance', (req, res) => {
	    mReply.service.route = "maintenance";
    	res.json(mReply)
    });
	
	sApp.get("/testGet", (req, res) => {
		sReply.data.route = "testGet";
		return res.json(sReply);
	});
	
	sApp.post("/testPost", (req, res) => {
		sReply.data.route = "testPost";
		return res.json(sReply);
	});
	
	sApp.put("/testPut", (req, res) => {
		sReply.data.route = "testPut";
		return res.json(sReply);
	});
	
	sApp.delete("/testDelete", (req, res) => {
		sReply.data.route = "testPut";
		return res.json(sReply);
	});
	
	sApp.put("/testPut", (req, res) => {
		sReply.data.route = "testPut";
		return res.json(sReply);
	});
    sApp.patch("/testPatch", (req, res) => {
	    sReply.data.route = "testPatch";
        return res.json(sReply);
    });

    sApp.head("/testHead", (req, res) => {
	    sReply.data.route = "testHead";
	    return res.json(sReply);
    });

    sApp.options("/testOther", (req, res) => {
	    sReply.data.route = "testOther";
	    return res.json(sReply);
    });

    let sAppServer = sApp.listen(serverConfig.s.port, () => console.log(`${serverConfig.name} service listening on port ${serverConfig.s.port}!`));
    let mAppServer = mApp.listen(serverConfig.m.port, () => console.log(`${serverConfig.name} service listening on port ${serverConfig.m.port}!`));

    return callback(
        {
            "sAppServer": sAppServer,
            "mAppServer": mAppServer,
            "name": serverConfig.name
        }
    )
}

function stopServer(config) {
    console.log("Stopping server");

    config.mAppServer.close((err) => {
        console.log("...sAppServer: " + config.name);
    });

    config.sAppServer.close((err) => {
        console.log("...mAppServer: " + config.name);
    });
}

service.init(function () {
    let reg = service.registry.get();

    let dbConfig = reg.coreDB.provision;
    if (reg.coreDB.oauth) {
        dbConfig = {
            "provision": reg.coreDB.provision,
            "oauth": reg.coreDB.oauth
        };
    }
    provision.init(dbConfig, service.log);
});

startServer({s: {port: 4010}, m: {port: 5010}, name: "restApiService"}, function (servers) {

});