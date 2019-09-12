#lang racket

;The expression:

;f(g(x),h(y))
;is written in ANF as:

;let v0 = g(x) in
;    let v1 = h(y) in
;        f(v0,v1)


(define (anf-atom? expr)
  (match expr
    [x #:when (or (number? x) (symbol? x)) #t]
    [`(anf-lambda (,x) ,body) #t]
    [default #f]))


(define (anf expr)  
  (letrec ([id-gen (lambda (x) x)]
           [cps1 (lambda (expr codegen)
                   (match expr
                     [atom #:when (anf-atom? atom) (codegen expr)]
                     [`(,app ,rand) #:when (and (anf-atom? app) (anf-atom? rand))
                                    (if (eq? codegen id-gen)
                                        `(,app ,rand)
                                        (let ([x (gensym "var")])
                                          `(let ([,x (,app ,rand)])
                                             ,(codegen x))
                                          #;`(,app ,rand (lambda (,x) ,(codegen x)))))]
                     [`(,app ,rand) (cps1 rand (lambda (r-atom)
                                                 (cps1 app (lambda (a-atom)
                                                             (cps1 `(,a-atom ,r-atom) codegen)))))]
                     [`(lambda (,x) ,body) (codegen `(anf-lambda (,x) ,(anf body)))]))])
    (cps1 expr id-gen)))

(anf 'x)
(anf '(f x))
(anf '(x (y z)))
(anf '((m n) (f (o p))))
(anf '(lambda (x) (m x)))
(anf '((lambda (x) (m x)) (lambda (q) (lambda (a) (a (c b))))))
