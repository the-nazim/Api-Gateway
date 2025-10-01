from fastapi import FastAPI
from pydantic import BaseModel

app = FastAPI()

class Product(BaseModel):
    id: int 
    name: str

@app.get("/product/{product_id}", response_model=Product)
async def getProduct(product_id: int):
    if product_id == 1:
        return Product(id=1, name="Sample Product")
    return {"id": product_id, "name": "Unknown Product"}