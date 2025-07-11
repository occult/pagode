import { useState, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { CreditCardIcon, LockIcon } from "lucide-react";

interface PaymentFormProps {
  plan: {
    id: string;
    name: string;
    price: number;
    currency: string;
  };
  onSubmit: (paymentMethodId: string) => void;
  isProcessing: boolean;
  stripePublishableKey: string;
}

export function PaymentForm({ plan, onSubmit, isProcessing, stripePublishableKey }: PaymentFormProps) {
  const [stripe, setStripe] = useState<any>(null);
  const [elements, setElements] = useState<any>(null);
  const [cardElement, setCardElement] = useState<any>(null);
  const [error, setError] = useState<string>("");
  const [name, setName] = useState("");

  const formatPrice = (price: number, currency: string) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: currency.toUpperCase(),
    }).format(price / 100);
  };

  useEffect(() => {
    // Initialize Stripe
    const initializeStripe = async () => {
      if (window.Stripe) {
        const stripeInstance = window.Stripe(stripePublishableKey);
        setStripe(stripeInstance);

        const elementsInstance = stripeInstance.elements();
        setElements(elementsInstance);

        const cardElementInstance = elementsInstance.create('card', {
          style: {
            base: {
              fontSize: '16px',
              color: '#424770',
              '::placeholder': {
                color: '#aab7c4',
              },
            },
          },
        });

        cardElementInstance.mount('#card-element');
        setCardElement(cardElementInstance);

        cardElementInstance.on('change', (event: any) => {
          setError(event.error ? event.error.message : '');
        });
      }
    };

    // Load Stripe script if not already loaded
    if (!window.Stripe) {
      const script = document.createElement('script');
      script.src = 'https://js.stripe.com/v3/';
      script.onload = initializeStripe;
      document.head.appendChild(script);
    } else {
      initializeStripe();
    }
  }, [stripePublishableKey]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!stripe || !cardElement || !name.trim()) {
      return;
    }

    setError('');

    // Create payment method
    const { error: paymentMethodError, paymentMethod } = await stripe.createPaymentMethod({
      type: 'card',
      card: cardElement,
      billing_details: {
        name: name,
      },
    });

    if (paymentMethodError) {
      setError(paymentMethodError.message);
      return;
    }

    // Pass the payment method ID to the parent component
    onSubmit(paymentMethod.id);
  };

  return (
    <Card className="w-full max-w-md mx-auto">
      <CardHeader className="text-center">
        <div className="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-primary/10">
          <CreditCardIcon className="h-6 w-6 text-primary" />
        </div>
        <CardTitle>Payment Details</CardTitle>
        <CardDescription>
          Subscribe to {plan.name} for {formatPrice(plan.price, plan.currency)}/month
        </CardDescription>
      </CardHeader>
      
      <CardContent>
        <form onSubmit={handleSubmit} className="space-y-4">
          {/* Cardholder Name */}
          <div className="space-y-2">
            <label htmlFor="name" className="text-sm font-medium">Cardholder Name</label>
            <input
              id="name"
              type="text"
              placeholder="John Doe"
              value={name}
              onChange={(e) => setName(e.target.value)}
              required
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>

          {/* Card Element */}
          <div className="space-y-2">
            <label className="text-sm font-medium">Card Information</label>
            <div 
              id="card-element" 
              className="p-3 border border-gray-300 rounded-md focus-within:ring-2 focus-within:ring-blue-500"
            />
            {error && (
              <div className="text-sm text-red-600">{error}</div>
            )}
          </div>

          {/* Test Card Notice */}
          <div className="p-3 bg-blue-50 border border-blue-200 rounded-md dark:bg-blue-950 dark:border-blue-800">
            <div className="flex items-start gap-2">
              <LockIcon className="h-4 w-4 text-blue-600 mt-0.5 flex-shrink-0" />
              <div className="text-sm text-blue-800 dark:text-blue-200">
                <p className="font-medium">Test Mode</p>
                <p className="text-xs mt-1">
                  Use test card: 4242 4242 4242 4242, any future expiry date, any CVC.
                </p>
              </div>
            </div>
          </div>

          {/* Submit Button */}
          <Button 
            type="submit" 
            className="w-full mt-6" 
            disabled={isProcessing || !stripe || !name.trim()}
          >
            {isProcessing ? (
              <div className="flex items-center gap-2">
                <div className="h-4 w-4 animate-spin rounded-full border-2 border-current border-t-transparent" />
                Processing Payment...
              </div>
            ) : (
              `Subscribe for ${formatPrice(plan.price, plan.currency)}/month`
            )}
          </Button>

          {/* Security Notice */}
          <div className="text-xs text-muted-foreground text-center">
            <LockIcon className="h-3 w-3 inline mr-1" />
            Your payment information is secure and encrypted
          </div>
        </form>
      </CardContent>
    </Card>
  );
}