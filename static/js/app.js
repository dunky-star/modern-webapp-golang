// Prompt is our JavaScript module for all alerts, notifications, and custom popup dialogs
function Prompt() {
    let toast = function (c) {
        const {
            msg = '',
            icon = 'success',
            position = 'top-end',

        } = c;

        const Toast = Swal.mixin({
            toast: true,
            title: msg,
            position: position,
            icon: icon,
            showConfirmButton: false,
            timer: 3000,
            timerProgressBar: true,
            didOpen: (toast) => {
                toast.addEventListener('mouseenter', Swal.stopTimer)
                toast.addEventListener('mouseleave', Swal.resumeTimer)
            }
        })

        Toast.fire({})
    }

    let success = function (c) {
        const {
            msg = "",
            title = "",
            footer = "",
        } = c

        Swal.fire({
            icon: 'success',
            title: title,
            text: msg,
            footer: footer,
        })

    }

    let error = function (c) {
        const {
            msg = "",
            title = "",
            footer = "",
        } = c

        Swal.fire({
            icon: 'error',
            title: title,
            text: msg,
            footer: footer,
        })

    }

    async function custom(c) {
        const {
            icon = "",
            msg = "",
            title = "",
            showConfirmButton = true,
        } = c;

        const {value: result} = await Swal.fire({
            icon: icon,
            title: title,
            html: msg,
            backdrop: false,
            focusConfirm: false,
            showCancelButton: true,
            showConfirmButton: showConfirmButton,
            willOpen: () => {
                if (c.willOpen !== undefined) {
                    c.willOpen();
                }
            },
            didOpen: () => {
                if (c.didOpen !== undefined) {
                    c.didOpen();
                }
            },
            preConfirm: () => {
                return [
                    document.getElementById('start').value,
                    document.getElementById('end').value
                ]
            }
        })

        if (result) {
            if (result.dismiss !== Swal.DismissReason.cancel) {
                if (result.value !== "") {
                    if (c.callback !== undefined) {
                        c.callback(result);
                    }
                } else {
                    c.callback(false);
                }
            } else {
                c.callback(false);
            }
        }
    }

    return {
        toast: toast,
        success: success,
        error: error,
        custom: custom,
    }
}

// initRoomCheckAvailability wires up the "Check Availability" button on room pages.
// roomId: room id (e.g. 1 or 2). csrfToken: CSRF token from the template (e.g. "{{.CSRFToken}}").
function initRoomCheckAvailability(roomId, csrfToken) {
    var button = document.getElementById("check-availability-button");
    if (!button) return;

    var formHtml = [
        '<form id="check-availability-form" action="" method="post" novalidate class="needs-validation">',
        '  <div class="form-row"><div class="col">',
        '    <div class="form-row" id="reservation-dates-modal">',
        '      <div class="col"><input disabled required class="form-control" type="text" name="start" id="start" placeholder="Arrival"></div>',
        '      <div class="col"><input disabled required class="form-control" type="text" name="end" id="end" placeholder="Departure"></div>',
        '    </div>',
        '  </div></div>',
        '</form>'
    ].join('');

    button.addEventListener("click", function () {
        attention.custom({
            msg: formHtml,
            title: "Choose your dates",
            willOpen: function () {
                var elem = document.getElementById("reservation-dates-modal");
                if (elem) {
                    new DateRangePicker(elem, { format: "yyyy-mm-dd", showOnFocus: true, minDate: new Date() });
                }
            },
            didOpen: function () {
                var start = document.getElementById("start");
                var end = document.getElementById("end");
                if (start) start.removeAttribute("disabled");
                if (end) end.removeAttribute("disabled");
            },
            callback: function (result) {
                if (!result || !Array.isArray(result) || result.length < 2) return;
                var formData = new FormData();
                formData.append("start", result[0] || "");
                formData.append("end", result[1] || "");
                formData.append("csrf_token", csrfToken);
                formData.append("room_id", String(roomId));

                fetch("/search-availability-json", {
                    method: "post",
                    body: formData,
                    credentials: "same-origin"
                })
                    .then(function (r) { return r.json(); })
                    .then(function (data) {
                        if (data.ok) {
                            attention.custom({
                                icon: "success",
                                showConfirmButton: false,
                                msg: '<p>Room is available!</p><p><a href="/book-room?id=' + data.room_id + '&s=' + data.start_date + '&e=' + data.end_date + '" class="btn btn-primary">Book now!</a></p>'
                            });
                        } else {
                            attention.error({ msg: "No availability" });
                        }
                    });
            }
        });
    });
}
