(load "mk.scm")
(import (Framework match))
(import (Framework helpers))
(case-sensitive #t)
(print-gensym 'pretty)

(define-syntax λ
  (syntax-rules ()
    [(_ (x ...) body ...)
     (lambda (x ...) body ...)]))


(define (id x) x)

(define (k0 x) `(k ,x))

(define (cps p k)
  (match p
    [,sym (guard (symbol? sym)) (k sym)]
    [(λ (,x) ,body) `(λ (,x k) ,(cps body k0))]
    [(add1 ,expr) (cps expr `(K-add1 ,k)
                       (λ (ecode)
                              (k `(add1 ,ecode))))]
    [(,app ,rator)
     (cps app (λ (acode)
                (cps rator (λ (rcode)
                             `(,acode ,rcode
                                      ,(if (eq? k k0)
                                           'k
                                           (let ([sym (gensym)])
                                             `(λ (,sym)
                                                ,(k sym)))))))))]))



(define (cps p k)
  (match p
    [,sym (guard (symbol? sym)) (apply-k k sym)]
    [(λ (,x) ,body) `(λ (,x k) ,(cps body 'K-k0))]
    [(add1 ,expr) (cps expr (λ (ecode)
                              (apply-k k `(add1 ,ecode))))]
    [(,app ,rator)
     (cps app `(K-app-1 ,k ,rator))]))


(define (apply-k k code)
  (match k
    [K-id code]
    [K-k0 `(k ,code)]
    [(K-add1 ,k) (apply-k k `(add1 ,code))]
    [(K-app-1 ,k ,rator)
     (cps rator `(K-app-0 ,k ,code))]
    [(K-app-0 ,k ,acode)
     `(,acode ,code
              ,(if (eq? k 'K-k0)
                   'k
                   (let ([sym (gensym)])
                     `(λ (,sym)
                        ,(apply-k k sym)))))]))




;;;;;
(define (cpsᵒ p k out)
  (conde
   [(symbolo p) (apply-kᵒ k p out)]
   [(fresh (x body bout)
      (== p `(λ (,x) ,body))
      (== `(λ (,x k) ,bout) out)
      (cpsᵒ body 'K-k0 bout))]
   [(fresh (app rator)
      (== p `(,app ,rator))
      (cpsᵒ app `(K-app-1 ,k ,rator) out))]))


(define (apply-kᵒ k-strcut code out)
  (conde
   [(== k-strcut 'K-id) (== code out)]
   [(== k-strcut 'K-k0) (== `(k ,code) out)]
   [(fresh (k)
      (== k-strcut `(K-add1 ,k))
      (apply-kᵒ k `(add1 ,code) out))]
   [(fresh (k rator)
      (== k-strcut `(K-app-1 ,k ,rator))
      (cpsᵒ rator `(K-app-0 ,k ,code) out))]
   [(fresh (k acode)
      (== k-strcut `(K-app-0 ,k ,acode))
      (conde
       [(== k 'K-k0) (== `(,acode ,code k) out)]
       [(fresh (sym appout)
          (symbolo sym)
          (apply-kᵒ k sym appout)
          (== `(,acode ,code (λ (,sym) ,appout)) out))]))]))



#!eof
(load "the-cpser.ss")

(cps `(p (λ (x) ((z p) q))) id)

(cps 'x id)

(run 1 (res k)
    (cpsᵒ res k 'x))
