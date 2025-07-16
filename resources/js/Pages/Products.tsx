import AppLayout from "@/Layouts/AppLayout";
import { type BreadcrumbItem } from "@/types";
import { Head, router } from "@inertiajs/react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { useState, useEffect } from "react";
import { PaymentForm } from "@/components/PaymentForm";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { CheckCircle2 } from "lucide-react";

const breadcrumbs: BreadcrumbItem[] = [
  {
    title: "Dashboard",
    href: "/dashboard",
  },
  {
    title: "Products",
    href: "/products",
  },
];

interface Product {
  id: string;
  name: string;
  description: string;
  price: number;
  currency: string;
}

interface ProductsProps {
  title: string;
  product: Product;
  stripePublishableKey: string;
  success?: boolean;
  paymentIntentId?: string;
}


export default function Products({ title, product, stripePublishableKey, success, paymentIntentId }: ProductsProps) {
  const [showPaymentModal, setShowPaymentModal] = useState(false);
  const [paymentError, setPaymentError] = useState<string>("");
  const [showSuccess, setShowSuccess] = useState(false);
  const [processing, setProcessing] = useState(false);

  const formatPrice = (price: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
    }).format(price / 100);
  };


  const handlePaymentSubmit = (paymentMethodId: string) => {
    setProcessing(true);
    
    router.post('/products/purchase', {
      productId: product.id,
      paymentMethodId: paymentMethodId,
      amount: product.price,
    }, {
      onSuccess: () => {
        setShowPaymentModal(false);
        setShowSuccess(true);
      },
      onError: (errors) => {
        setPaymentError(errors.message || "Payment failed. Please try again.");
      },
      onFinish: () => {
        setProcessing(false);
      },
    });
  };

  // Check for success prop
  useEffect(() => {
    if (success) {
      setShowSuccess(true);
    }
  }, [success]);

  return (
    <AppLayout breadcrumbs={breadcrumbs}>
      <Head title={title} />
      <div className="flex h-full flex-1 flex-col gap-4 rounded-xl p-4">
        <h1 className="text-3xl font-bold">Products</h1>
        
        {showSuccess && (
          <Alert className="border-green-200 bg-green-50">
            <CheckCircle2 className="h-4 w-4 text-green-600" />
            <AlertDescription className="text-green-800">
              Payment successful! Thank you for your purchase.
            </AlertDescription>
          </Alert>
        )}
        
        <div className="grid gap-6 md:grid-cols-1 lg:grid-cols-2 xl:grid-cols-3">
          <Card className="w-full">
            <CardHeader>
              <CardTitle>{product.name}</CardTitle>
              <CardDescription>{product.description}</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="text-3xl font-bold">{formatPrice(product.price)}</div>
              <p className="text-sm text-muted-foreground mt-2">
                One-time purchase
              </p>
            </CardContent>
            <CardFooter>
              <Button 
                onClick={() => setShowPaymentModal(true)}
                className="w-full"
                size="lg"
              >
                Purchase Now
              </Button>
            </CardFooter>
          </Card>
        </div>

        <Dialog open={showPaymentModal} onOpenChange={setShowPaymentModal}>
          <DialogContent className="sm:max-w-md">
            <DialogHeader>
              <DialogTitle>Complete Purchase</DialogTitle>
            </DialogHeader>
            <div className="space-y-4">
              <div className="flex justify-between items-center py-2 border-b">
                <span className="font-medium">{product.name}</span>
                <span className="font-bold">{formatPrice(product.price)}</span>
              </div>
              
              <PaymentForm
                plan={{
                  id: product.id,
                  name: product.name,
                  price: product.price,
                  currency: product.currency,
                }}
                onSubmit={handlePaymentSubmit}
                isProcessing={processing}
                stripePublishableKey={stripePublishableKey}
                mode="payment"
              />
              
              {paymentError && (
                <p className="text-sm text-red-600">{paymentError}</p>
              )}
              
              <div className="flex gap-2 pt-4">
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => setShowPaymentModal(false)}
                  className="flex-1"
                  disabled={processing}
                >
                  Cancel
                </Button>
              </div>
            </div>
          </DialogContent>
        </Dialog>
      </div>
    </AppLayout>
  );
}