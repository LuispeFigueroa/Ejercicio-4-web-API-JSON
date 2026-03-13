Ejercicio 4, API JSON Luis Pedro Figueroa #24087

- EL api esta basado en luchadores de la UFC, especificamente en la division de peso ligero, cada luchador tiene id,nombre,record,especialidad y nacionalidad.
- LA API implementa 4 operaciones: Get,Post,Patch y Delete. 
- Todos los datos se almacenan en un archivo json en la direccion "data/fighters.json"
- El servidor se ejecuta en el puerto de mi carnet, que es 24087.

- Ejemplo de la estructura de los luchadores
{
  "id": 1,
  "name": "Conor McGregor",
  "country": "Ireland",
  "record": "22-6-0",
  "specialty": "Striking",
  "height": "1.75m"
}

Endpoint:

- Obtener a todos los luchadores de la lista:
GET http://localhost:24087/api/Fighters

- Get con query parameters la busqueda en este caso no es case-sensitive:
GET /api/Fighters?country=Brazil
GET /api/Fighters?specialty=Grappling
GET /api/Fighters?country=Brazil&specialty=Grappling

- Crear un nuevo fighter:
POST /api/Fighters



- actualizar datos de un luchador
PATCH /api/Fighters?id=1

-Eliminar a un luchador
DELETE /api/Fighters&id=1

-Todos los errores se devuelven en formato JSON estructurado:
{
  "error": "Invalid ID"
}