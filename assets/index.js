const refreshFrequencyMilliseconds = 2000;

let isSet = undefined;

async function render() {
    const response = await fetch('http://localhost:8090/status');
    const data = await response.json();

    function createRowHtml({ id, name, pid, running, auto_restart, group }) {
        
        const runningTextColor = running ? 'text-success' : 'text-danger';
        const runningTextStatus = running ? 'Running' : 'Stopped';

        return `
            <div class="row py-2">
                <div class="col-2">${name}</div>
                <div class="col-2">${pid}</div>
                <div class="col-2 ${runningTextColor} fw-bold">${runningTextStatus}</div>
                <div class="col-1">${auto_restart}</div>
                <div class="col-2">${group}</div>
                <div class="col-3">
                    <button onclick="run('${id}')" class="btn btn-success">Start</button>
                    <button onclick="stop('${id}')" class="btn btn-danger">Stop</button>
                    <button onclick="showLogsModal('${id}')" class="btn btn-info">Logs</button>
                </div>
            </div>
        `;
    }

    function createGroupHtml(items) {

        const group = items[0].group;
        return `
            <div class="row py-2">
                <div class="col-12 text-center">
                    <hr>
                </div>
            </div>
            <div class="row py-2">
                <div class="col-2 fw-bold">
                    Group With Id = ${group}
                </div>
                <div class="col-2">
                    <button onclick="rungroup('${group}')" class="btn btn-success">Start Group</button>
                </div>
                <div class="col-2">
                    <button onclick="stopgroup('${group}')" class="btn btn-danger">Stop Group</button>
                </div>
            </div>
            ${items.map(createRowHtml).join('')}
        `;
    }

    isSet = data.length != 0;

    const header = `
        <div class="row py-2 mb-5">
            <div class="col-2">
                <button onclick="runall()" class="btn btn-success">Start All</button>
            </div>
            <div class="col-2">
                <button onclick="stopall()" class="btn btn-danger">Stop All</button>
            </div>
            <div class="col-2">
                <button onclick="toggleSet()" class="btn btn-warning">${isSet ? 'Unset' : 'Set'}</button>
            </div>
        </div>
        <div class="row py-2">
            <div class="col-2">Name</div>
            <div class="col-2">PID</div>
            <div class="col-2">Running</div>
            <div class="col-1">Auto Restart</div>
            <div class="col-2">Group</div>
            <div class="col-3">Actions</div>
        </div>
    `;

    // --- Split data into groups ---

    // for example: { 'group X': 0, 'group Y': 1, 'group alphabetical order doesn't matter': 2 }
    const groups = {}; 

    // for example: 
    // [  [ { group: 'group X', ... }, { group: 'group X', ... } ], 
    //    [ { group: 'group Y' } ], 
    //    [ { group: 'group alphabetical order doesn't matter' } ]  ]
    const groupsList = [];

    for (const item of data) {
        if (!(item.group in groups)) {
            groups[item.group] = groupsList.length;
            groupsList.push([]);
        }
        groupsList[groups[item.group]].push(item);
    }

    // --- Now groupsList contains the items grouped by group preserving their order ---

    const html = [header, ...groupsList.map(createGroupHtml)].join('');
    
    $('#main-content').html(html);
}

async function showLogsModal(id) {

    const logsRequest = getAllLogs(id);
    $('#logsModalBody').html('Loading....');
    $('#logsModal').modal('show');
    const logs = await logsRequest;
    function renderLogLine(line) {
        const lineColor = line.type === 'errors' ? 'text-danger' : '';
        return `<span class="${lineColor}">${line.message}</span>`;
    }
    const html = logs.map(renderLogLine).join('<br>');
    $('#logsModalBody').html(`<p>${html}</p>`);
}

let triggerRefreshImmediately = () => {};

async function main() {

    while (true) {
        try {
            await render();
        }
        catch (error) {
            console.error('Error:', error);
        }

        await new Promise(resolve => {
            triggerRefreshImmediately = resolve;
            setTimeout(resolve, refreshFrequencyMilliseconds);
        });
        
        triggerRefreshImmediately = () => {};
    }
}

function alertAsync(message) {
    setTimeout(() => alert(message), 100);
}

async function simpleRequest(request) {

    try {
        const response = await fetch(request);
        triggerRefreshImmediately();
        const responseJson = await response.json();
        alertAsync(responseJson.message);
    }
    catch (error) {
        alertAsync('Error executing the request, see console for more details');
        console.error('Error:', error);
    }
}

async function run(id) {
    await simpleRequest(`/run?id=${id}`);
}

async function stop(id) {
    await simpleRequest(`/stop?id=${id}`);
}

async function runGroup(group) {
    await simpleRequest(`/rungroup?group=${group}`);
}

async function stopGroup(group) {
    await simpleRequest(`/stopgroup?group=${group}`);
}

async function runAll() {
    await simpleRequest(`/runall`);
}

async function stopAll() {
    await simpleRequest(`/stopall`);
}

async function toggleSet() {
    if (isSet === false) await simpleRequest(`/set`);
    if (isSet === true) await simpleRequest(`/unset`);
}


/**
 * 
 * @param {string} id 
 * @param {"out" | "errors"} type 
 * @param {Number} offset 
 * @returns 
 */
async function getLogs(id, type, offset) {
    
    const response = await fetch(`/execlogs?id=${id}&type=${type}&offset=${offset}`);
    const responseText = await response.text();
    if (!response.ok) {
        throw new Error(`Failed to get logs for id=${id}, type=${type}, offset=${offset}, response: ${responseText}`);
    }
    return responseText.trim().split('\n').map(line => ({ timestamp: '', type: type, message: line }));
}

async function getAllLogs(id) {

    let logs = [];
    try {
        while (true) {

            const moreLogs = await getLogs(id, 'errors', logs.length);

            logs = [...moreLogs, ...logs];
        }
    }
    catch (error) {
        // maybe handle/throw errors other than "offset out of range" here
    }
    return logs;
}


main();
