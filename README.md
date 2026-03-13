Ejercicio 4, API JSON Luis Pedro Figueroa #24087

- EL api esta basado en luchadores de la UFC, especificamente en la division de peso ligero, cada luchador tiene id,nombre,record,especialidad y nacionalidad.
- LA API implementa 4 operaciones: Get,Post,Patch y Delete. 
- Todos los datos se almacenan en un archivo json en la direccion "data/fighters.json"
- El servidor se ejecuta en el puerto de mi carnet, que es 24087.
  <img width="657" height="188" alt="evidencia docker corriendo en el puerto de mi carnet " src="https://github.com/user-attachments/assets/bfc6c6ff-e833-4c28-ab41-5f9c02ef5954" />

- Ejemplo de la estructura de los luchadores

{
  "id": 1,
  "name": "Conor McGregor",
  "country": "Ireland",
  "record": "22-6-0",
  "specialty": "Striking",
  "height": "1.75m"
}

Endpoints:



- Obtener a todos los luchadores de la lista:
GET http://localhost:24087/api/Fighters

<img width="1918" height="1079" alt="evidencia get " src="https://github.com/user-attachments/assets/fb34180e-674a-4216-ab5b-1dbd725f9ca5" />




- Get con query parameters la busqueda en este caso no es case-sensitive:
GET /api/Fighters?country=Brazil,
GET /api/Fighters?specialty=Grappling,
GET /api/Fighters?country=Brazil&specialty=Grappling,

<img width="1915" height="1036" alt="evidencia Get con otro parametro" src="https://github.com/user-attachments/assets/6a0af052-722c-4c05-a029-ce4d97286cb7" />



-GET con path parameters:
GET /api/Fighters/1

<img width="1894" height="1027" alt="evidencia path parameters " src="https://github.com/user-attachments/assets/f11023e5-6970-4cf7-9a89-3f8fde225044" />




- Crear un nuevo fighter:
POST /api/Fighters

<img width="1918" height="1079" alt="captura post" src="https://github.com/user-attachments/assets/514b9517-4c02-4076-b633-980ea68afb68" />





- actualizar datos de un luchador
PATCH /api/Fighters?id=1

<img width="1911" height="1037" alt="evidencia patch" src="https://github.com/user-attachments/assets/4a2daed5-5ae0-4753-99db-ec326a03453a" />

-Eliminar a un luchador
DELETE /api/Fighters&id=1

<img width="1908" height="1029" alt="evidencia delete " src="https://github.com/user-attachments/assets/0ebd3d9e-6b69-470a-828b-03dbd65466cd" />



-Todos los errores se devuelven en formato JSON estructurado:
{
  "error": "Invalid ID"
}

<img width="1915" height="1027" alt="Evidencia de caso de error" src="https://github.com/user-attachments/assets/6ec363ec-5e32-449b-8ba5-df4c305066bd" />


