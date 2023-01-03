(def fib
  ((fn rec-fib [a b]
     (lazy-seq (cons a (rec-fib b (+ a b)))))
   0 1))
