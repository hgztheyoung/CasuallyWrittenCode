#lang racket

#;(x -> (k x))
#;(((m n) -> (m n k)))
#;((m n) -> (m n (lambda (mnv) (k mnv))))
#;((x (y z)) -> (y z (lambda (yzres) (x yzres k))))
#;(((lambda (x) (m x)) (lambda (a) b)) -> ((lambda (x k0) (m x k0)) (k0 (lambda (a k1) (k1 b))) k))

(define (atom-call? call)
  (member call `(add1 sub1)))

(define (atom? expr)
  (match expr
    [x #:when (or (number? x) (symbol? x)) #t]
    [`(lambda (,x k) ,body) #t]
    [default #f]))


(define (cps expr)  
  (letrec ([id-gen (lambda (expr) `(k ,expr))]
           [cps1 (lambda (expr codegen)
                   (match expr
                     [atom #:when (atom? atom) (codegen expr)]
                     [`(,app ,rand) #:when (and (atom? app) (atom? rand))
                                     (if (eq? codegen id-gen)
                                         `(,app ,rand k) 
                                         (let ([x (gensym "var")])
                                           `(,app ,rand (lambda (,x) ,(codegen x)))))]
                     [`(,app ,rand) (cps1 rand (lambda (r-atom)
                                                   (cps1 app (lambda (a-atom)
                                                               (cps1 `(,a-atom ,r-atom) codegen)))))]
                     [`(lambda (,x) ,body) (codegen `(lambda (,x k) ,(cps body)))]))])
    (cps1 expr id-gen)))

(cps 'x)
(cps '(f x))
(cps '(x (y z)))
(cps '((m n) (o p)))
(cps '(lambda (x) (m x)))
(cps '((lambda (x) (m x)) (lambda (q) (lambda (a) b))))



(define (cps-without-eta-optimize expr)  
  (letrec ([id-gen (lambda (expr) `(k ,expr))]
           [cps1 (lambda (expr codegen)
                   (match expr
                     [atom #:when (atom? atom) (codegen expr)]
                     [`(,app ,rand) #:when (and (atom? app) (atom? rand))
                                     (let ([x (gensym "var")])
                                       `(,app ,rand (lambda (,x) ,(codegen x))))]
                     [`(,app ,rand) (cps1 rand (lambda (rpiece)
                                                   (cps1 app (lambda (apiece)
                                                               (cps1 `(,apiece ,rpiece) codegen)))))]
                     [`(lambda (,x) ,body) (codegen `(lambda (,x k) ,(cps1 body id-gen)))]))])
    (cps1 expr id-gen)))

(cps-without-eta-optimize '(m n))
