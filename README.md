# Внешняя сортировка слиянием
_Для сортировки фрагментов используется метод сортировки выбором_

## Допущения:
1. Файл содержит исключительно печатные символы из таблицы ASCII
2. Считаем, что порядок символом должен соответствовать порядку их расположения в таблице ASCII
3. В условиях временных ограничений не ставилось целью сделать хорошее оформление кода и разбиение на структуры
4. Решение не претендует быть оптимальным и требует доработки (как архитектурной, так и логической)

## Выбранный метод решения
1. Исходный файл разбивается на фрагменты по 1000 строк
2. Каждый фрагмент сортируется методом сортировки выбора и сохраняется во временную папку
3. Отсортированные фрагменты объединяются методом сортировки слиянием до тех пор, пока не останется одного файла

## Как пользоваться
1. Установите [Go](https://golang.org/dl/)
2. ...