export const handleNewShortURL = async (url) => {
    const response = await fetch(import.meta.env.VITE_API_URL, {
        method: "PUT",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            url,
        }),
    }).then((res) => res.json());

    const error = response.status !== 201 ? true : false;
    const messages = response.errors
        ? response.errors.map((res) => res.message)
        : response.message;

    return [response, error, messages]
};