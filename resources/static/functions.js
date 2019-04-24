// Namespace.
var pushgateway = {};

pushgateway.labels = {};
pushgateway.panel = null;

pushgateway.switchToMetrics = function () {
    $('#metrics-div').removeClass('hidden');
    $('#status-div').addClass('hidden');
    $('#metrics-li').addClass('active');
    $('#status-li').removeClass('active');
};

pushgateway.switchToStatus = function () {
    $('#metrics-div').addClass('hidden');
    $('#status-div').removeClass('hidden');
    $('#metrics-li').removeClass('active');
    $('#status-li').addClass('active');
};

pushgateway.showDelModal = function (labels, panelID, event) {
    event.stopPropagation(); // Don't trigger accordion collapse.
    pushgateway.labels = labels;
    pushgateway.panel = $('#' + panelID);

    var components = [];
    for (var ln in labels) {
        components.push(ln + '="' + labels[ln] + '"')
    }

    $('#del-modal-msg').text(
        'Do you really want to delete all metrics of group {' + components.join(', ') + '}?'
    );
    $('#del-modal').modal('show');
};

pushgateway.deleteGroup = function () {
    var pathElements = [];
    for (var ln in pushgateway.labels) {
        if (ln != 'job') {
            pathElements.push(encodeURIComponent(ln));
            pathElements.push(encodeURIComponent(pushgateway.labels[ln]));
        }
    }
    var groupPath = pathElements.join('/');
    if (groupPath.length > 0) {
        groupPath = '/' + groupPath;
    }

    $.ajax({
        type: 'DELETE',
        url: 'metrics/job/' + encodeURIComponent(pushgateway.labels['job']) + groupPath,
        success: function (data, textStatus, jqXHR) {
            pushgateway.panel.remove();
            $('#del-modal').modal('hide');

            // Disable "del all" button if we just deleted the last job row
            let jobs = $("#job-accordion").children("div");
            if (jobs.length === 0) {
                $('#del-all-button').addClass("disabled");
            }
        },
        error: function (jqXHR, textStatus, error) {
            alert('Deleting metric group failed: ' + error);
        }
    });
};

pushgateway.showDelAllModal = function () {
    $('#del-all-modal').modal('show');
};

pushgateway.deleteAllGroups = function () {
    // Retrieve all job panels
    const $panels = $("#job-accordion>.panel>.panel-heading");

    const allRawLabels = [];
    const promises = [];

    // For each job panel
    for (let i = 0; i < $panels.length; i++) {
        const $panel = $($panels[i]);
        const $h4 = $panel.children("h4.panel-title");

        // Get the labels content (including job)
        const $spans = $h4.children("span.label");

        const jobRawLabels = [];

        // For each label span, extract the label name and the label value
        for (let j = 0; j < $spans.length; j++) {
            const rawLabel = $spans[j].innerText;
            jobRawLabels.push(rawLabel);
        }

        // Add the job raw labels to the global list
        allRawLabels.push(jobRawLabels);

        // Path building, encoding URI component and removing label value's quotes
        const pathElements = ["metrics"];
        jobRawLabels.forEach((kv) => {
            let [key, value] = kv.split("=");
            value = value.replace(/"/g, "");

            if (key != "job") {
                key = encodeURIComponent(key);
                value = encodeURIComponent(value);
            }

            pathElements.push(key, value)
        });

        // Creation of the actual path string
        let path = pathElements.join("/");

        // Sending request and store it into a list of promises
        promises.push(fetch(path, {method: "DELETE"}));
    }

    // Wait for of all promises to finish, or at least one fail
    Promise.all(promises)
        .then(() => {
            // hide the mode, remove the job rows and disable the "del all" button
            $('#del-all-modal').modal('hide');
            $('#job-accordion')[0].innerHTML = "";
            $('#del-all-button').addClass('disabled');
        })
        .catch((error) => {
            // Display an error message and do nothing else
            alert('Deleting all metric groups failed: ' + error);
        })
};
