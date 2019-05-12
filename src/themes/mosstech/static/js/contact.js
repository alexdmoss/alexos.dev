jQuery(document).ready(function (e) {
    /* Contacts Form */
    $(function () {
        $('#contacts').find('input,select,textarea').jqBootstrapValidation({
            preventSubmit: true,
            submitError: function ($form, event, errors) {
            },
            submitSuccess: function ($form, e) {
                e.preventDefault()
                var submitButton = $('input[type=submit]', $form)
                // if you're copying this code, the API below still won't work even though you have the URL + API Key ;)
                $.ajax({
                    type: 'POST',
                    url: 'https://contact.alexos.dev/api/email/mosstech.io',
                    headers: {
                        'API-Key': 'H3mWAM.ouZZHVtDKPj9iD3pVSjmTlabq70lcJmv'
                    },
                    data: $form.serialize(),
                    beforeSend: function (xhr, opts) {
                        if ($('#_email', $form).val()) {
                            xhr.abort()
                        } else {
                            submitButton.prop('value', 'Please Wait...')
                            submitButton.prop('disabled', 'disabled')
                        }
                    }
                }).done(function (data) {
                    submitButton.prop('value', 'Thanks for your message!')
                    submitButton.prop('disabled', true)
                })
            },
            filter: function () {
                return $(this).is(':visible')
            }
        })
    })
})
$('#name').focus(function () {
    $('#success').html('')
})