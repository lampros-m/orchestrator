const refreshFrequencyMilliseconds = 2000;

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
                <div class="col-2">${auto_restart}</div>
                <div class="col-2">${group}</div>
                <div class="col-2">
                    <button onclick="run('${id}')" class="btn btn-success">Start</button>
                    <button onclick="stop('${id}')" class="btn btn-danger">Stop</button>
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

    const header = `
        <div class="row py-2 mb-5">
            <div class="col-2">
                <button onclick="runall()" class="btn btn-success">Start All</button>
            </div>
            <div class="col-2">
                <button onclick="stopall()" class="btn btn-danger">Stop All</button>
            </div>
        </div>
        <div class="row py-2">
            <div class="col-2">Name</div>
            <div class="col-2">PID</div>
            <div class="col-2">Running</div>
            <div class="col-2">Auto Restart</div>
            <div class="col-2">Group</div>
            <div class="col-2">Actions</div>
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

    console.log(groupsList);

    // --- Now groupsList contains the items grouped by group preserving their order ---

    const html = [header, ...groupsList.map(createGroupHtml)].join('');
    
    $('#main-content').html(html);
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

async function simpleRequest(request) {

    try {
        const response = await fetch(request);
        triggerRefreshImmediately();
        const responseJson = await response.json();
        alert(responseJson.message);
    }
    catch (error) {
        alert('Error executing the request, see console for more details');
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


main();