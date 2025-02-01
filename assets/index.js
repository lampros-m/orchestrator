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

    const header = `
        <div class="row py-2">
            <div class="col-2">Name</div>
            <div class="col-2">PID</div>
            <div class="col-2">Running</div>
            <div class="col-2">Auto Restart</div>
            <div class="col-2">Group</div>
            <div class="col-2">Actions</div>
        </div>
    `;

    const html = [header, ...data.map(createRowHtml)].join('');
    
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

main();