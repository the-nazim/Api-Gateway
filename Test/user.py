from fastapi import FastAPI
from pydantic import BaseModel

app = FastAPI()
class User(BaseModel):
    id: int 
    name: str

@app.get("/user/{user_id}", response_model=User)
async def getUser(user_id: int):
    if user_id == 1:
        return User(id=1, name="John Doe")
    return {"id": user_id, "name": "Unknown"}