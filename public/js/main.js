copy = (txt) => {
    navigator.clipboard.writeText(txt);
}
copyShort = (txt) => {
    copy(window.location.origin + "/" + txt);
}
formCtrl = () => {
    fetch("/api", {
        method: "POST",
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            short_path: document.getElementById('form-shorten').value,
            real_url: document.getElementById('form-real').value,
        })
    }).then(
        response => response.text()
    ).then((html) => {
        console.log(html);
        window.location.reload();
    }
    );
}
deleteShort = (short) => {
    fetch(`/api/${short}`, {
        method: "DELETE",
    }).then(
        response => response.text()
    ).then(html => {
        html => console.log(html);
        window.location.reload();
    });
}