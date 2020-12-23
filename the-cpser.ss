(import (Framework match))
(import (Framework helpers))
(case-sensitive #t)
(optimize-level 2)
(print-gensym 'pretty)
(load "mk.scm")


(define-syntax λ
  (syntax-rules ()
    [(_ (x ...) body ...)
     (lambda (x ...) body ...)]))


(define (cps1 program k)
  (match program
    [,s (guard (symbol? s)) (k s)]
    [(λ (,x) ,body) (guard (symbol? x))
                    (k `(λ (,x k) ,(cps1 body (λ (x) `(k ,x)))))]
    [(,app ,rator)
     (cps1 app (λ (acode)
                 (cps1 rator (λ (rcode)
                               (let ([sym (gensym)])
                                 `(,acode ,rcode (λ (,sym)
                                                   ,(k sym))))))))]))

(define (subst o n s-exp)
  (match s-exp
    [() '()]
    [,a (guard (symbol? a)) (if (eq? a o) n a)]
    [(,[a] . ,[d]) `(,a . ,d)]))


(define (substᵒ o n s-exp out)
  (conde
   [(== s-exp '()) (== '() out)]
   [(symbolo s-exp)
    (conde
     [(== s-exp o) (== n out)]
     [(=/= s-exp o) (== s-exp out)])]
   [(fresh (a d ares dres)
      (== `(,ares . ,dres) out)
      (== `(,a . ,d) s-exp)
      (substᵒ o n a ares)
      (substᵒ o n d dres))]))

(define (cps-i cpsed-program k)
  (match cpsed-program
    [,s (guard (symbol? s)) (k s)]
    [(λ (,x k) ,body) (guard (symbol? x))
                      (cps-i body (λ (rbody)
                                    (k `(λ (,x) ,rbody))))]
    [(k ,sth) (k (cps-i sth k))]
    [(,app ,rator (λ (,ressymbol) ,body))
     (cps-i app (λ (rapp)
                  (cps-i rator (λ (rrator)
                                 (let* ([body-res (cps-i body k)])
                                   (subst ressymbol `(,rapp ,rrator) body-res))))))]))

; defunc cps-i
;we treat subst as atomic here to save some strength
(define (cps-i-defunc cpsed-program k)
  (match cpsed-program
    [,s (guard (symbol? s)) (apply-cps-i-defunc k s)]
    [(λ (,x k) ,body) (guard (symbol? x))
                      (cps-i-defunc body `(K-λ-final ,x ,k))]
    [(k ,sth) (cps-i-defunc sth k)]
    [(,app ,rator (λ (,sym) ,body))
     (cps-i-defunc app
                   `(K-app-rator ,rator ,sym ,body ,k))]))

(define (apply-cps-i-defunc k-struct code)
  (match k-struct
    [K-id code]
    [(K-λ-final ,x ,k) (apply-cps-i-defunc k `(λ (,x) ,code))]
    [(K-app-rator ,rator ,sym ,body ,k)
     (cps-i-defunc rator `(K-app-body ,code ,sym ,body ,k))]
    [(K-app-body ,rapp ,sym ,body ,k)
     (cps-i-defunc body `(K-app-final ,rapp ,code ,sym ,k))]
    [(K-app-final ,rapp ,rrator ,sym ,k)
     (apply-cps-i-defunc k (subst sym `(,rapp ,rrator) code))]))


(define-syntax display-expr
  (syntax-rules ()
    [(_ expr) (begin
                (display expr)
                (display "\n")
                expr)]))


(define (decent-λ p bound-l)
  (match p
    [,sym (guard (member sym bound-l)) sym]
    [(λ (,x) ,body) (guard (symbol? x))
                    `(λ (,x) ,(decent-λ body (cons x bound-l)))]
    [(,[app] ,[rator])
     `(,app ,rator)]))



(define (consᵒ a d p)
  (== p `(,a . ,d)))


(define (carᵒ l a)
  (fresh (d)
    (consᵒ a d l)))


(define (cdrᵒ l d)
  (fresh (a)
    (consᵒ a d l)))


(define (pairᵒ p)
  (fresh (a d)
    (consᵒ a d p)))

(define (listᵒ l)
  (conde
   [(== l '())]
   [(pairᵒ l)
    (fresh (d)
      (cdrᵒ l d)
      (listᵒ d))]))


(define (proper-memberᵒ x l)
  (conde
   [(carᵒ l x)
    (fresh (d)
      (cdrᵒ l d)
      (listᵒ d))]
   [(fresh (d)
      (cdrᵒ l d)
      (proper-memberᵒ x d))]))

(define (decent-λᵒ bound-syms out)
  (conde
   [(symbolo out)
    (proper-memberᵒ out bound-syms)]
   [(fresh (x body nm)
      (== out `(λ (,x) ,body))
      (symbolo x)
      (consᵒ x bound-syms nm)
      (decent-λᵒ nm body))]
   [(fresh (app rator)
      (== out `(,app ,rator))
      (decent-λᵒ bound-syms app)
      (decent-λᵒ bound-syms rator))]))

(define (decent-λ p bound-l)
  (match p
    [,sym (guard (member sym bound-l)) sym]
    [(λ (,x) ,body) (guard (symbol? x))
                    `(λ (,x) ,(decent-λ body (cons x bound-l)))]
    [(,[app] ,[rator])
     `(,app ,rator)]))

(define (cpsᵒ p k out)
  (conde
   [(symbolo p) (apply-kᵒ k p out)]
   [(fresh (x body bout)
      (== p `(λ (,x) ,body))
      (cpsᵒ body 'K-k0 bout)
      (apply-kᵒ k `(λ (,x k) ,bout) out))]
   [(fresh (app rator)
      (== p `(,app ,rator))
      (cpsᵒ app `(K-app-1 ,k ,rator) out))]))


(define (apply-kᵒ k-strcut code out)
  (conde
   [(== k-strcut 'K-id) (== code out)]
   [(== k-strcut 'K-k0) (== `(k ,code) out)]
   [(fresh (k rator)
      (== k-strcut `(K-app-1 ,k ,rator))
      (cpsᵒ rator `(K-app-0 ,k ,code) out))]
   [(fresh (k acode)
      (== k-strcut `(K-app-0 ,k ,acode))
      (conda
       [(== k 'K-k0) (== `(,acode ,code k) out)]
       [(fresh (sym appout)
          (symbolo sym)
          (apply-kᵒ k sym appout)
          (== `(,acode ,code (λ (,sym) ,appout)) out))]))]))

(define (decent-cpsedᵒ p)
  (fresh (raw)
    (decent-λᵒ '() raw)
    (cpsᵒ raw 'K-id p)))

(define (id x) x)

;now,ready to create cps-iᵒ
(define (cps-iᵒ program k out)
  (decent-cpsedᵒ program)
  (conde
   [(symbolo program) (apply-cps-iᵒ k program out)]
   [(fresh (x body)
      (== program `(λ (,x k) ,body))
      (symbolo x)
      (cps-iᵒ body `(K-λ-final ,x ,k) out))]
   [(fresh (sth)
      (== program `(k ,sth))
      (cps-iᵒ sth k out))]
   [(fresh (app rator sym body)
      (== program `(,app ,rator (λ (,sym) ,body)))
      (cps-iᵒ app `(K-app-rator ,rator ,sym ,body ,k) out))]))

(define (apply-cps-iᵒ k-struct code out)
  (conde
   [(== k-struct 'K-id) (== code out)]
   [(fresh (x k)
      (== k-struct `(K-λ-final ,x ,k))
      (apply-cps-iᵒ k `(λ (,x) ,code) out))]
   [(fresh (rator sym body k)
      (== k-struct `(K-app-rator ,rator ,sym ,body ,k))
      (cps-iᵒ rator `(K-app-body ,code ,sym ,body ,k) out))]
   [(fresh (rapp sym body k)
      (== k-struct `(K-app-body ,rapp ,sym ,body ,k))
      (cps-iᵒ body `(K-app-final ,rapp ,code ,sym ,k) out))]
   [(fresh (rapp rrator sym k subst-res)
      (== k-struct `(K-app-final ,rapp ,rrator ,sym ,k))
      (substᵒ sym `(,rapp ,rrator) code subst-res)
      (apply-cps-iᵒ k subst-res out))]))


#!eof

(load "the-cpser.ss")

(run 40 (res)
  (decent-cpsedᵒ res))


(define bound-λs
 (map car
       (run 400 (res)
         (decent-λᵒ res '() res))))

(map (λ (p) (cps1 p id))
     bound-λs)

(trace decent-λ)

(decent-λ '((λ (x) (λ (y) x)) (λ (z) z)) '())

(cps1 '(λ (x) (λ (y) x)) id)


(run 40 (res)
  (cps-iᵒ res 'K-id 'x))


(run* (out)
      (cps-iᵒ '(p (λ (x k) (k z)) (λ (g0) (g0 q (λ (g1) g1))))
              'K-id out))


(run 1 (out)
  (cps-iᵒ '(p (λ (x k) (k z)) (λ (g0) (g0 q (λ (g1) g1))))
          out
          '(((p (λ (x) z)) q))))




(cps-i '(λ (x k) (k x)) (λ (x) x))


(cps-i
 '(p (λ (x k) (k z)) (λ (g4) (g4 q (λ (g5) g5))))
 (λ (x) x))

(cps-i-defunc
 '(p (λ (x k) (k z)) (λ (g4) (g4 q (λ (g5) g5))))
 'K-id)

(run* (q)
      (substᵒ 'x 'y '(x y z) q))

(run 1 (out)
  (cps-iᵒ
   '(p (λ (x k) (k z)) (λ (g4) (g4 q (λ (g5) g5))))
   'K-id
   out))


(run 10 (out)
      (cps-iᵒ
       out
       'K-id
      '(λ (x) x)))


(run* (out)
      (cps-iᵒ
       '(λ (x k) (k x))
       'K-id
       out))

(run* (out)
      (cps-iᵒ
       '(p q (λ (x) x))
       'K-id
       out))

(run* (out)
      (cps-iᵒ '(p (λ (x k) (k z)) (λ (g0) (g0 q (λ (g1) g1))))
              'K-id out))



(cps1 `((p (λ (x) z)) q) (λ (x) x))

(define-syntax display-expr
  (syntax-rules ()
    [(_ expr) (begin
                (display expr)
                (display "\n")
                                expr)]))
