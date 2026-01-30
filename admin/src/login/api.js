export async function login(username, password, rememberMe) {
    const res = await fetch('/admin/api/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password, remember_me: rememberMe }),
    });
    return res.json();
}
