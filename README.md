# spendtracker

A golang program to transform bank statements to a CSV report that aggregates expenses by date and tags.

Here is an example. Given a statement file like this:

```
01/12/2017,"-300.00",Tetsuya Restaruant
15/11/2017,"-100.00",WOOLWORTHS
15/12/2017,"-100.00",WOOLWORTHS
16/12/2017,"10.00",WOOLWORTHS REFUND
31/12/2017,"-10.00",UNITED - Visa Purchase - Receipt xxxx Card xxx
25/01/2018,"-10.00",UNITED - Visa Purchase - Receipt xxxx Card xxx
```

And a pattern file matching regular expressions to tags:

```
Woolworths,Living Expenses,Groceries
Tetsuya,Non-essential Expenses,Restaurant
UNITED ,Living Expenses,Fuel
```

It will generate a report like this:

```
2017-11,2017-12,2018-01
Living Expenses,Groceries,-100.00,-90.00,0.00
Living Expenses,Fuel,0.00,-10.00,-20.00
Non-essential Expenses,Restaurant,0.00,-300.00,0.00
```

From here this report can be imported into a spreadsheet and further analysed.

