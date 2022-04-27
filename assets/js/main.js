// GET POST PUT DELETE
function Send(method, url, data, callback) {
    let xhr = new XMLHttpRequest();
    xhr.open(method, url);
    xhr.onload = function (event) {
        callback && callback(JSON.parse(this.response));
    };
    xhr.setRequestHeader("Content-Type", "application/json; charset=utf-8");
    xhr.send(JSON.stringify(data));
}
// UPLOAD FILE TO BACKEND
function Upload(file, callback) {
    let xhr = new XMLHttpRequest();
    xhr.open("POST", "/upload");
    let data = new FormData();
    data.append("MyFiles", file, file.name);
    xhr.onload = function (event) {
        callback && callback(JSON.parse(this.response));
    }
    xhr.send(data);
}

// Добавление нового "продукта"
function CreateProduct() {
    let file = ProductImage.files[0];
    // Если файл изображения не был выбран, то ничего не делаем
    if (!file) return;

    // Загрузить файл на сервер
    Upload(file, ()=> {
        // После загрузки отправить запрос на добавление
        // самого товара в базу
        Send("PUT", "/product", {
            Name: ProductName.value,
            Price: +ProductPrice.value,
            Image: "/image/" + file.name
        }, () => {
            // Т.к. в базе было пополнение товаров
            // то вызвать функцию загрузки товаров
            // чтобы обновить список
            GetProducts('категория1');
        });
    });
}

function GetProducts(category) {
    Send("GET", "/product/"+category, null, (response) => {
        List.innerHTML = "";
        for (let item of response) {
            let divItem = document.createElement("div");
            divItem.className = "item";

            let image = document.createElement("img");
            image.src = item.Image;

            let h2 = document.createElement("h2");
            h2.textContent = item.Name;

            let p = document.createElement("p");
            p.textContent = ""+item.Price + " рублей";

            divItem.append(image, h2, p);

            List.append(divItem);
        }
    });
}
GetProducts();

