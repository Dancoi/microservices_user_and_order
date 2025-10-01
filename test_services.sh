#!/bin/bash

echo "=== ТЕСТИРОВАНИЕ USER SERVICE ==="

# 1. Создаем первого пользователя
echo "1. Создаем пользователя Alice:"
curl -X POST http://localhost:3001/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice"}' \
  -w "\nStatus: %{http_code}\n\n"

sleep 1

# 2. Создаем второго пользователя  
echo "2. Создаем пользователя Bob:"
curl -X POST http://localhost:3001/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Bob"}' \
  -w "\nStatus: %{http_code}\n\n"

sleep 1

# 3. Получаем всех пользователей
echo "3. Все пользователи:"
curl http://localhost:3001/users \
  -w "\nStatus: %{http_code}\n\n"

sleep 1

# 4. Получаем пользователя по ID=1
echo "4. Пользователь с ID=1:"
curl http://localhost:3001/users/1 \
  -w "\nStatus: %{http_code}\n\n"

sleep 1

# 5. Получаем пользователя по ID=2
echo "5. Пользователь с ID=2:"
curl http://localhost:3001/users/2 \
  -w "\nStatus: %{http_code}\n\n"

sleep 2

echo "=== ТЕСТИРОВАНИЕ ORDER SERVICE ==="

# 6. Создаем первый заказ
echo "6. Создаем заказ для Alice (user_id=1):"
curl -X POST http://localhost:3002/orders \
  -H "Content-Type: application/json" \
  -d '{"userId":1, "item":"MacBook Pro"}' \
  -w "\nStatus: %{http_code}\n\n"

sleep 1

# 7. Создаем второй заказ
echo "7. Создаем заказ для Bob (user_id=2):"
curl -X POST http://localhost:3002/orders \
  -H "Content-Type: application/json" \
  -d '{"userId":2, "item":"iPhone 15"}' \
  -w "\nStatus: %{http_code}\n\n"

sleep 1

# 8. Получаем все заказы
echo "8. Все заказы:"
curl http://localhost:3002/orders \
  -w "\nStatus: %{http_code}\n\n"

sleep 1

# 9. Получаем заказ по ID=1
echo "9. Заказ с ID=1 (с данными пользователя):"
curl http://localhost:3002/orders/1 \
  -w "\nStatus: %{http_code}\n\n"

sleep 1

# 10. Получаем заказ по ID=2
echo "10. Заказ с ID=2 (с данными пользователя):"
curl http://localhost:3002/orders/2 \
  -w "\nStatus: %{http_code}\n\n"

echo "=== ТЕСТИРОВАНИЕ ЗАВЕРШЕНО ==="