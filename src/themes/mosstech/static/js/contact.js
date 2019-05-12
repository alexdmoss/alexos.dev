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
                $.ajax({
                    type: 'POST',
                    url: '{{ .Site.Params.contactAPI }}',
                    headers: {
                        'Origin': '{{ .Site.Title }}',
                        'API-Key': '{{ .Site.Params.contactKey }}'
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